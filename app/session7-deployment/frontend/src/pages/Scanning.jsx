import { useState, useEffect } from "react";
import {
  Play,
  RefreshCw,
  Clock,
  CheckCircle,
  XCircle,
  AlertCircle,
  Activity,
  Eye,
} from "lucide-react";
import { assetsAPI, scanningAPI } from "../services/api";

const SCAN_TYPES = [
  { value: "all", label: "All Passive Scans", color: "primary", passive: true },
  { value: "dns", label: "DNS Records", color: "info", passive: true },
  { value: "whois", label: "WHOIS Lookup", color: "success", passive: true },
  { value: "subdomain", label: "Subdomain Enumeration", color: "warning", passive: true },
  { value: "cert_trans", label: "Certificate Transparency", color: "info", passive: true },
  { value: "asn", label: "ASN Lookup", color: "primary", passive: true },
  { value: "ip", label: "IP Geolocation", color: "info", passive: true },
  { value: "port", label: "⚠️ Port Scan (Active)", color: "danger", passive: false },
  { value: "ssl", label: "⚠️ SSL/TLS Probe (Active)", color: "danger", passive: false },
];

function Scanning() {
  const [assets, setAssets] = useState([]);
  const [selectedAsset, setSelectedAsset] = useState("");
  const [selectedScanType, setSelectedScanType] = useState("dns");
  const [scanJobs, setScanJobs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [scanning, setScanning] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  // Results modal state
  const [selectedJob, setSelectedJob] = useState(null);
  const [jobResults, setJobResults] = useState(null);
  const [loadingResults, setLoadingResults] = useState(false);
  const [showModal, setShowModal] = useState(false);

  useEffect(() => {
    loadAssets();
  }, []);

  useEffect(() => {
    if (selectedAsset) {
      loadScanJobs();
      const interval = setInterval(loadScanJobs, 5000);
      return () => clearInterval(interval);
    }
  }, [selectedAsset]);

  const loadAssets = async () => {
    try {
      const data = await assetsAPI.list({ status: "active", page_size: 100 });
      setAssets(data.data || []);
      if (data.data && data.data.length > 0) {
        setSelectedAsset(data.data[0].id);
      }
    } catch (err) {
      setError(err.message);
    }
  };

  const loadScanJobs = async () => {
    if (!selectedAsset) return;
    try {
      setLoading(true);
      const data = await scanningAPI.listJobs(selectedAsset, { page: 1, page_size: 20 });
      setScanJobs(data.data || []);
    } catch (err) {
      console.error("Failed to load scan jobs:", err);
    } finally {
      setLoading(false);
    }
  };

  const handleStartScan = async () => {
    if (!selectedAsset || !selectedScanType) return;

    const scanType = SCAN_TYPES.find((t) => t.value === selectedScanType);
    if (!scanType.passive) {
      const confirmed = window.confirm(
        "⚠️ WARNING: You are about to start an ACTIVE scan.\n\n" +
          "Active scans directly probe target systems and may be illegal without authorization.\n\n" +
          "Only proceed if you own the target or have written permission.\n\n" +
          "Continue?"
      );
      if (!confirmed) return;
    }

    try {
      setScanning(true);
      setError("");
      setSuccess("");
      await scanningAPI.startScan(selectedAsset, selectedScanType);
      setSuccess("Scan started successfully! Refreshing...");
      setTimeout(() => {
        loadScanJobs();
        setSuccess("");
      }, 2000);
    } catch (err) {
      setError(err.message);
    } finally {
      setScanning(false);
    }
  };

  const handleViewResults = async (job) => {
    try {
      setLoadingResults(true);
      setSelectedJob(job);
      setShowModal(true);
      const data = await scanningAPI.getResults(job.id);
      setJobResults(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoadingResults(false);
    }
  };

  const closeModal = () => {
    setShowModal(false);
    setSelectedJob(null);
    setJobResults(null);
  };

  const getStatusIcon = (status) => {
    switch (status) {
      case "pending": return <Clock size={16} className="status-pending" />;
      case "running": return <RefreshCw size={16} className="status-running" />;
      case "completed": return <CheckCircle size={16} className="status-completed" />;
      case "failed": return <XCircle size={16} className="status-failed" />;
      case "partial": return <AlertCircle size={16} className="status-warning" />;
      default: return null;
    }
  };

  const getStatusBadge = (status) => {
    const badges = {
      pending: "badge-warning",
      running: "badge-info",
      completed: "badge-success",
      failed: "badge-danger",
      partial: "badge-warning",
    };
    return badges[status] || "badge-secondary";
  };

  const selectedAssetData = assets.find((a) => a.id === selectedAsset);

  // ============================================
  // RENDER RESULTS IN MODAL
  // ============================================

  const renderSSLResults = (results) => {
    if (!results?.results?.length) return <p className="text-muted">No SSL data available.</p>;
    return results.results.map((r, i) => (
      <div key={i}>
        <div className="flex items-center gap-2 mb-3">
          <span className={`badge ${r.grade === "A" ? "badge-success" : r.grade === "B" ? "badge-warning" : "badge-danger"}`}>
            Grade: {r.grade}
          </span>
          <strong>{r.domain}</strong>
        </div>
        <table className="table">
          <tbody>
            <tr><td><strong>Subject</strong></td><td>{r.certificate?.subject}</td></tr>
            <tr><td><strong>Issuer</strong></td><td>{r.certificate?.issuer}</td></tr>
            <tr><td><strong>Valid From</strong></td><td>{r.certificate?.valid_from ? new Date(r.certificate.valid_from).toLocaleDateString() : "N/A"}</td></tr>
            <tr><td><strong>Valid Until</strong></td><td>{r.certificate?.valid_until ? new Date(r.certificate.valid_until).toLocaleDateString() : "N/A"}</td></tr>
            <tr><td><strong>Days Until Expiry</strong></td><td>{r.certificate?.days_until_expiry} days</td></tr>
            <tr><td><strong>Is Expired</strong></td><td>{r.certificate?.is_expired ? "Yes ❌" : "No ✅"}</td></tr>
            <tr><td><strong>Self Signed</strong></td><td>{r.certificate?.is_self_signed ? "Yes" : "No"}</td></tr>
            <tr><td><strong>TLS Version</strong></td><td>{r.connection?.tls_version}</td></tr>
            <tr><td><strong>Cipher Suite</strong></td><td>{r.connection?.cipher_suite}</td></tr>
            <tr><td><strong>SANs</strong></td><td style={{wordBreak:"break-all"}}>{r.certificate?.san?.join(", ")}</td></tr>
          </tbody>
        </table>
        {r.issues?.length > 0 && (
          <div className="alert alert-error mt-2">⚠️ Issues: {r.issues.join(", ")}</div>
        )}
      </div>
    ));
  };

  const renderPortResults = (results) => {
    if (!results?.results?.length) return <p className="text-muted">No port scan data available.</p>;
    return results.results.map((r, i) => (
      <div key={i}>
        <div className="mb-3 p-3" style={{background:"var(--color-background-secondary)", borderRadius:"8px"}}>
          <p className="text-sm"><strong>IP:</strong> {r.ip_address}</p>
          <p className="text-sm"><strong>Total Scanned:</strong> {r.total_scanned} ports</p>
          <p className="text-sm"><strong>Open Ports:</strong> {r.open_ports?.length || 0}</p>
          <p className="text-sm"><strong>Closed Ports:</strong> {r.closed_ports}</p>
          <p className="text-sm"><strong>Duration:</strong> {r.scan_duration_ms}ms</p>
        </div>
        {!r.open_ports?.length ? (
          <p className="text-muted">No open ports found.</p>
        ) : (
          <table className="table">
            <thead>
              <tr><th>Port</th><th>Protocol</th><th>State</th><th>Service</th></tr>
            </thead>
            <tbody>
              {r.open_ports.map((p, j) => (
                <tr key={j}>
                  <td><strong>{p.port}</strong></td>
                  <td>{p.protocol}</td>
                  <td><span className="badge badge-success">{p.state}</span></td>
                  <td>{p.service}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    ));
  };

  const renderIPResults = (results) => {
    if (!results?.results?.length) return <p className="text-muted">No IP geolocation data available.</p>;
    return results.results.map((r, i) => (
      <div key={i}>
        <table className="table">
          <tbody>
            <tr><td><strong>IP Address</strong></td><td>{r.ip_address}</td></tr>
            <tr><td><strong>Country</strong></td><td>{r.geolocation?.country} ({r.geolocation?.country_code})</td></tr>
            <tr><td><strong>City / Region</strong></td><td>{r.geolocation?.city}, {r.geolocation?.region}</td></tr>
            <tr><td><strong>ISP</strong></td><td>{r.geolocation?.isp}</td></tr>
            <tr><td><strong>Organization</strong></td><td>{r.geolocation?.org}</td></tr>
            <tr><td><strong>ASN Number</strong></td><td>{r.asn?.number}</td></tr>
            <tr><td><strong>ASN Name</strong></td><td>{r.asn?.name}</td></tr>
            <tr><td><strong>Reverse DNS</strong></td><td>{r.reverse_dns || "N/A"}</td></tr>
            <tr><td><strong>Coordinates</strong></td><td>{r.geolocation?.latitude}, {r.geolocation?.longitude}</td></tr>
          </tbody>
        </table>
      </div>
    ));
  };

  const renderDNSResults = (results) => {
    const records = Array.isArray(results) ? results : results?.data || [];
    if (!records.length) return <p className="text-muted">No DNS records found.</p>;
    return (
      <table className="table">
        <thead><tr><th>Type</th><th>Name</th><th>Value</th><th>TTL</th></tr></thead>
        <tbody>
          {records.map((r, i) => (
            <tr key={i}>
              <td><span className="badge badge-primary">{r.record_type}</span></td>
              <td>{r.name}</td>
              <td className="text-sm" style={{wordBreak:"break-all"}}>{r.value}</td>
              <td className="text-sm text-muted">{r.ttl}s</td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  };

  const renderModalContent = () => {
    if (loadingResults) {
      return <div className="loading"><div className="spinner"></div><span>Loading results...</span></div>;
    }
    if (!jobResults) return <p className="text-muted">No results available.</p>;

    switch (selectedJob?.scan_type) {
      case "ssl": return renderSSLResults(jobResults);
      case "port": return renderPortResults(jobResults);
      case "ip": return renderIPResults(jobResults);
      case "dns": return renderDNSResults(jobResults);
      default:
        return (
          <pre style={{fontSize:"12px", overflow:"auto", maxHeight:"400px"}}>
            {JSON.stringify(jobResults, null, 2)}
          </pre>
        );
    }
  };

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">Scanning Operations</h1>
        <p className="page-description">Execute reconnaissance scans on your assets</p>
      </div>

      {error && <div className="alert alert-error mb-4">{error}</div>}
      {success && <div className="alert alert-success mb-4">{success}</div>}

      <div className="grid grid-2">
        {/* Scan Configuration */}
        <div className="card">
          <div className="card-header">
            <h3 className="card-title">Start New Scan</h3>
          </div>

          {assets.length === 0 ? (
            <div className="empty-state">
              <p className="text-muted">No active assets found. Please add assets first.</p>
            </div>
          ) : (
            <>
              <div className="form-group">
                <label className="form-label">Select Asset</label>
                <select className="form-select" value={selectedAsset} onChange={(e) => setSelectedAsset(e.target.value)}>
                  {assets.map((asset) => (
                    <option key={asset.id} value={asset.id}>{asset.name} ({asset.type})</option>
                  ))}
                </select>
              </div>

              <div className="form-group">
                <label className="form-label">Scan Type</label>
                <select className="form-select" value={selectedScanType} onChange={(e) => setSelectedScanType(e.target.value)}>
                  {SCAN_TYPES.map((type) => (
                    <option key={type.value} value={type.value}>
                      {type.label} {type.passive ? "🟢" : "🔴"}
                    </option>
                  ))}
                </select>
              </div>

              {selectedScanType && !SCAN_TYPES.find((t) => t.value === selectedScanType)?.passive && (
                <div className="alert alert-warning mb-4">
                  <strong>⚠️ Active Scan Warning:</strong> This scan type directly probes the target and requires explicit authorization.
                  Only proceed if you own the target or have written permission.
                </div>
              )}

              <button className="btn btn-primary w-full" onClick={handleStartScan} disabled={scanning || !selectedAsset}>
                {scanning ? (
                  <><RefreshCw size={18} className="animate-spin" /> Starting...</>
                ) : (
                  <><Play size={18} /> Start Scan</>
                )}
              </button>

              {selectedAssetData && (
                <div className="mt-4 p-4 bg-gray-50 rounded-md">
                  <h4 className="font-semibold text-sm mb-2">Target Details:</h4>
                  <div className="text-sm text-muted space-y-1">
                    <p><strong>Name:</strong> {selectedAssetData.name}</p>
                    <p><strong>Type:</strong> {selectedAssetData.type}</p>
                    <p><strong>Status:</strong> {selectedAssetData.status}</p>
                  </div>
                </div>
              )}
            </>
          )}
        </div>

        {/* Scan Types Info */}
        <div className="card">
          <div className="card-header">
            <h3 className="card-title">Available Scan Types</h3>
          </div>
          <div className="space-y-3">
            <div>
              <h4 className="font-semibold text-sm mb-2">🟢 Passive Scans (Safe)</h4>
              <ul className="text-sm text-muted space-y-1" style={{ listStyle: "disc", paddingLeft: "1.5rem" }}>
                <li><strong>All:</strong> Run all passive scans</li>
                <li><strong>DNS:</strong> Query public DNS records</li>
                <li><strong>WHOIS:</strong> Lookup domain registration</li>
                <li><strong>Subdomain:</strong> Enumerate subdomains</li>
                <li><strong>Cert Trans:</strong> Certificate Transparency logs</li>
                <li><strong>ASN:</strong> Autonomous System lookup</li>
                <li><strong>IP:</strong> Geolocation & ASN for IP assets</li>
              </ul>
            </div>
            <div className="pt-3 border-t">
              <h4 className="font-semibold text-sm mb-2">🔴 Active Scans (Requires Permission)</h4>
              <ul className="text-sm text-muted space-y-1" style={{ listStyle: "disc", paddingLeft: "1.5rem" }}>
                <li><strong>Port:</strong> TCP/UDP port scanning (localhost only)</li>
                <li><strong>SSL:</strong> SSL/TLS certificate probing</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      {/* Recent Scan Jobs */}
      <div className="card mt-4">
        <div className="card-header flex items-center justify-between">
          <h3 className="card-title">Recent Scan Jobs</h3>
          {selectedAsset && (
            <button className="btn btn-sm btn-secondary" onClick={loadScanJobs} disabled={loading}>
              <RefreshCw size={14} className={loading ? "animate-spin" : ""} /> Refresh
            </button>
          )}
        </div>

        {loading && scanJobs.length === 0 ? (
          <div className="loading"><div className="spinner"></div><span>Loading scan jobs...</span></div>
        ) : scanJobs.length === 0 ? (
          <div className="empty-state">
            <Activity className="empty-state-icon" size={64} />
            <h3 className="empty-state-title">No scans yet</h3>
            <p className="empty-state-description">Start your first scan to see results here</p>
          </div>
        ) : (
          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>Scan Type</th>
                  <th>Status</th>
                  <th>Started</th>
                  <th>Duration</th>
                  <th>Results</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {scanJobs.map((job) => (
                  <tr key={job.id}>
                    <td>
                      <div className="flex items-center gap-2">
                        {getStatusIcon(job.status)}
                        <span className="font-medium">{job.scan_type}</span>
                      </div>
                    </td>
                    <td>
                      <span className={`badge ${getStatusBadge(job.status)}`}>{job.status}</span>
                    </td>
                    <td className="text-sm text-muted">
                      {new Date(job.started_at).toLocaleString()}
                    </td>
                    <td className="text-sm text-muted">
                      {job.ended_at
                        ? `${Math.round((new Date(job.ended_at) - new Date(job.started_at)) / 1000)}s`
                        : "-"}
                    </td>
                    <td>
                      <span className="font-semibold">{job.results}</span>
                    </td>
                    <td>
                      {job.status === "completed" && (
                        <button
                          className="btn btn-sm btn-secondary"
                          onClick={() => handleViewResults(job)}
                        >
                          <Eye size={14} /> View
                        </button>
                      )}
                      {job.status === "failed" && (
                        <span className="text-sm text-muted" title={job.error}>❌ Failed</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {/* Results Modal */}
      {showModal && selectedJob && (
        <div className="modal-overlay" onClick={closeModal}>
          <div
            className="modal"
            style={{ maxWidth: "750px", maxHeight: "85vh", overflowY: "auto" }}
            onClick={(e) => e.stopPropagation()}
          >
            <div className="modal-header">
              <h3 className="modal-title">
                {selectedJob.scan_type === "ssl" && "🔒 SSL/TLS"}
                {selectedJob.scan_type === "port" && "🔍 Port Scan"}
                {selectedJob.scan_type === "ip" && "🌍 IP Geolocation"}
                {selectedJob.scan_type === "dns" && "📡 DNS Records"}
                {!["ssl","port","ip","dns"].includes(selectedJob.scan_type) && selectedJob.scan_type.toUpperCase()}
                {" "}Results
              </h3>
              <button className="btn btn-sm btn-secondary" onClick={closeModal}>×</button>
            </div>
            <div className="modal-body">
              {renderModalContent()}
            </div>
            <div className="modal-footer">
              <button className="btn btn-secondary" onClick={closeModal}>Close</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default Scanning;
