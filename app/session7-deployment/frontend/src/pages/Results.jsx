import { useState, useEffect, useRef } from "react";
import {
  FileText,
  Globe,
  Server,
  Search,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { assetsAPI, resultsAPI, scanningAPI } from "../services/api";

function Results() {
  const [assets, setAssets] = useState([]);
  const [selectedAsset, setSelectedAsset] = useState("");
  const [resultType, setResultType] = useState("all");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // Separate state for each result type
  const [whoisData, setWhoisData] = useState(null);
  const [dnsData, setDnsData] = useState([]);
  const [subdomainData, setSubdomainData] = useState([]);

  // Scan job results (SSL, Port, IP)
  const [scanJobs, setScanJobs] = useState([]);
  const [selectedJobResults, setSelectedJobResults] = useState(null);
  const [selectedJob, setSelectedJob] = useState(null);
  const [loadingJobResults, setLoadingJobResults] = useState(false);

  // Pagination state for each type
  const [whoisPage, setWhoisPage] = useState(1);
  const [dnsPage, setDnsPage] = useState(1);
  const [subdomainPage, setSubdomainPage] = useState(1);

  // Filter/Search state for each type
  const [dnsSearch, setDnsSearch] = useState("");
  const [dnsTypeFilter, setDnsTypeFilter] = useState("");
  const [subdomainSearch, setSubdomainSearch] = useState("");
  const [subdomainActiveFilter, setSubdomainActiveFilter] = useState("");

  // Pagination metadata
  const [dnsPagination, setDnsPagination] = useState({});
  const [subdomainPagination, setSubdomainPagination] = useState({});

  const pageSize = 10;

  const isChangingAsset = useRef(false);

  useEffect(() => {
    loadAssets();
  }, []);

  useEffect(() => {
    if (selectedAsset) {
      isChangingAsset.current = true;
      setWhoisPage(1);
      setDnsPage(1);
      setSubdomainPage(1);
      setDnsSearch("");
      setDnsTypeFilter("");
      setSubdomainSearch("");
      setSubdomainActiveFilter("");
      setSelectedJob(null);
      setSelectedJobResults(null);

      if (resultType === "all") {
        loadAllResults(selectedAsset).finally(() => {
          setTimeout(() => { isChangingAsset.current = false; }, 100);
        });
      } else {
        loadSingleTypeResults().finally(() => {
          setTimeout(() => { isChangingAsset.current = false; }, 100);
        });
      }
    }
  }, [selectedAsset, resultType]);

  useEffect(() => {
    if (isChangingAsset.current) return;
    if (!selectedAsset || resultType !== "all") return;
    loadDNSResults();
  }, [dnsPage, dnsSearch, dnsTypeFilter]);

  useEffect(() => {
    if (isChangingAsset.current) return;
    if (!selectedAsset || resultType !== "all") return;
    loadSubdomainResults();
  }, [subdomainPage, subdomainSearch, subdomainActiveFilter]);

  useEffect(() => {
    if (selectedAsset && resultType === "dns") loadDNSResults();
  }, [selectedAsset, resultType, dnsPage, dnsSearch, dnsTypeFilter]);

  useEffect(() => {
    if (selectedAsset && resultType === "subdomains") loadSubdomainResults();
  }, [selectedAsset, resultType, subdomainPage, subdomainSearch, subdomainActiveFilter]);

  useEffect(() => {
    if (selectedAsset && resultType === "whois") loadWhoisResults();
  }, [selectedAsset, resultType]);

  useEffect(() => {
    if (selectedAsset && ["ssl", "port", "ip"].includes(resultType)) {
      loadScanJobs(selectedAsset);
    }
  }, [selectedAsset, resultType]);

  const loadAssets = async () => {
    try {
      const data = await assetsAPI.list({ page_size: 100 });
      setAssets(data.data || []);
      if (data.data && data.data.length > 0) {
        setSelectedAsset(data.data[0].id);
      }
    } catch (err) {
      setError(err.message);
    }
  };

  const loadAllResults = async (assetId = selectedAsset) => {
    setLoading(true);
    setError("");
    try {
      await loadWhoisResults(assetId);
      await loadDNSResults(assetId, { page: 1, search: "", typeFilter: "" });
      await loadSubdomainResults(assetId, { page: 1, search: "", activeFilter: "" });
      await loadScanJobs(assetId);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const loadScanJobs = async (assetId = selectedAsset) => {
    try {
      const data = await scanningAPI.listJobs(assetId, { page_size: 100 });
      const relevantJobs = (data.data || []).filter(
        (j) => j.status === "completed" && ["ssl", "port", "ip"].includes(j.scan_type)
      );
      setScanJobs(relevantJobs);
    } catch (err) {
      console.error("Failed to load scan jobs:", err);
    }
  };

  const handleViewJobResult = async (job) => {
    if (selectedJob?.id === job.id) {
      setSelectedJob(null);
      setSelectedJobResults(null);
      return;
    }
    try {
      setLoadingJobResults(true);
      setSelectedJob(job);
      const data = await scanningAPI.getResults(job.id);
      setSelectedJobResults(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoadingJobResults(false);
    }
  };

  const loadWhoisResults = async (assetId = selectedAsset) => {
    try {
      const data = await resultsAPI.getWHOIS(assetId);
      setWhoisData(data.data || data);
    } catch (err) {
      console.error("Failed to load WHOIS:", err);
      setWhoisData(null);
    }
  };

  const loadDNSResults = async (assetId = selectedAsset, options = {}) => {
    try {
      const { page = dnsPage, search = dnsSearch, typeFilter = dnsTypeFilter } = options;
      const params = { page, page_size: pageSize };
      if (search) params.search = search;
      if (typeFilter) params.type = typeFilter;
      const data = await resultsAPI.getDNS(assetId, params);
      setDnsData(data.data || []);
      setDnsPagination({
        total: data.total || 0,
        page: data.page || page,
        page_size: data.page_size || pageSize,
        total_pages: data.total_pages || 0,
      });
    } catch (err) {
      console.error("Failed to load DNS:", err);
      setDnsData([]);
      setDnsPagination({});
    }
  };

  const loadSubdomainResults = async (assetId = selectedAsset, options = {}) => {
    try {
      const { page = subdomainPage, search = subdomainSearch, activeFilter = subdomainActiveFilter } = options;
      const params = { page, page_size: pageSize };
      if (search) params.search = search;
      if (activeFilter !== "") params.active = activeFilter;
      const data = await resultsAPI.getSubdomains(assetId, params);
      setSubdomainData(data.data || []);
      setSubdomainPagination({
        total: data.total || 0,
        page: data.page || page,
        page_size: data.page_size || pageSize,
        total_pages: data.total_pages || 0,
      });
    } catch (err) {
      console.error("Failed to load Subdomains:", err);
      setSubdomainData([]);
      setSubdomainPagination({});
    }
  };

  const loadSingleTypeResults = async () => {
    if (!selectedAsset) return;
    try {
      setLoading(true);
      setError("");
      if (resultType === "dns") await loadDNSResults();
      else if (resultType === "subdomains") await loadSubdomainResults();
      else if (resultType === "whois") await loadWhoisResults();
      else if (["ssl", "port", "ip"].includes(resultType)) await loadScanJobs();
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const selectedAssetData = assets.find((a) => a.id === selectedAsset);

  // ============================================
  // RENDER HELPERS
  // ============================================

  const renderPagination = (pagination, currentPage, setPage) => {
    if (!pagination || !pagination.total_pages || pagination.total_pages <= 1) return null;
    const generatePageNumbers = () => {
      const pages = [];
      const totalPages = pagination.total_pages;
      if (totalPages <= 7) {
        for (let i = 1; i <= totalPages; i++) pages.push(i);
      } else {
        pages.push(1);
        if (currentPage > 3) pages.push("...");
        for (let i = Math.max(2, currentPage - 1); i <= Math.min(currentPage + 1, totalPages - 1); i++) pages.push(i);
        if (currentPage < totalPages - 2) pages.push("...");
        pages.push(totalPages);
      }
      return pages;
    };
    return (
      <div className="flex items-center justify-between mt-4 px-4 pb-4">
        <div className="text-sm text-muted">
          Showing {(currentPage - 1) * pageSize + 1} to{" "}
          {Math.min(currentPage * pageSize, pagination.total || 0)} of{" "}
          {pagination.total || 0} results
        </div>
        <div className="flex items-center gap-2">
          <button className="btn btn-secondary btn-sm" disabled={currentPage === 1} onClick={() => setPage(currentPage - 1)}>
            <ChevronLeft size={16} /> Previous
          </button>
          {generatePageNumbers().map((pageNum, idx) =>
            pageNum === "..." ? (
              <span key={`ellipsis-${idx}`} className="px-2 text-muted">...</span>
            ) : (
              <button key={pageNum} className={`btn btn-sm ${pageNum === currentPage ? "btn-primary" : "btn-secondary"}`} onClick={() => setPage(pageNum)}>
                {pageNum}
              </button>
            )
          )}
          <button className="btn btn-secondary btn-sm" disabled={currentPage === pagination.total_pages} onClick={() => setPage(currentPage + 1)}>
            Next <ChevronRight size={16} />
          </button>
        </div>
      </div>
    );
  };

  const renderDNSRecords = (records, showFilters = false) => (
    <>
      {showFilters && (
        <div className="p-4 border-b">
          <div className="grid grid-2 gap-4">
            <div>
              <label className="form-label">Search DNS Records</label>
              <div style={{ position: "relative" }}>
                <Search style={{ position: "absolute", left: "12px", top: "50%", transform: "translateY(-50%)", color: "var(--color-text-muted)", pointerEvents: "none", zIndex: 1 }} size={18} />
                <input type="text" className="form-input" style={{ paddingLeft: "40px" }} placeholder="Search by name or value..." value={dnsSearch}
                  onChange={(e) => { setDnsSearch(e.target.value); setDnsPage(1); }} />
              </div>
            </div>
            <div>
              <label className="form-label">Filter by Type</label>
              <select className="form-select" value={dnsTypeFilter} onChange={(e) => { setDnsTypeFilter(e.target.value); setDnsPage(1); }}>
                <option value="">All Types</option>
                <option value="A">A Records</option>
                <option value="AAAA">AAAA Records</option>
                <option value="MX">MX Records</option>
                <option value="NS">NS Records</option>
                <option value="TXT">TXT Records</option>
                <option value="CNAME">CNAME Records</option>
              </select>
            </div>
          </div>
        </div>
      )}
      {!records || records.length === 0 ? (
        <div className="p-4"><p className="text-muted">No DNS records found</p></div>
      ) : (
        <>
          <div className="table-container">
            <table className="table">
              <thead><tr><th>Type</th><th>Name</th><th>Value</th><th>TTL</th></tr></thead>
              <tbody>
                {records.map((record, idx) => (
                  <tr key={record.id || idx}>
                    <td><span className="badge badge-primary">{record.record_type}</span></td>
                    <td className="font-medium">{record.name}</td>
                    <td className="text-sm">{record.value}</td>
                    <td className="text-sm text-muted">{record.ttl}s</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {showFilters && renderPagination(dnsPagination, dnsPage, setDnsPage)}
        </>
      )}
    </>
  );

  const renderSubdomains = (subdomains, showFilters = false) => (
    <>
      {showFilters && (
        <div className="p-4 border-b">
          <div className="grid grid-2 gap-4">
            <div>
              <label className="form-label">Search Subdomains</label>
              <div style={{ position: "relative" }}>
                <Search style={{ position: "absolute", left: "12px", top: "50%", transform: "translateY(-50%)", color: "var(--color-text-muted)", pointerEvents: "none", zIndex: 1 }} size={18} />
                <input type="text" className="form-input" style={{ paddingLeft: "40px" }} placeholder="Search by subdomain name..." value={subdomainSearch}
                  onChange={(e) => { setSubdomainSearch(e.target.value); setSubdomainPage(1); }} />
              </div>
            </div>
            <div>
              <label className="form-label">Filter by Status</label>
              <select className="form-select" value={subdomainActiveFilter} onChange={(e) => { setSubdomainActiveFilter(e.target.value); setSubdomainPage(1); }}>
                <option value="">All Status</option>
                <option value="true">Active Only</option>
                <option value="false">Inactive Only</option>
              </select>
            </div>
          </div>
        </div>
      )}
      {!subdomains || subdomains.length === 0 ? (
        <div className="p-4"><p className="text-muted">No subdomains found</p></div>
      ) : (
        <>
          <div className="table-container">
            <table className="table">
              <thead><tr><th>Subdomain</th><th>Source</th><th>Active</th><th>Discovered</th></tr></thead>
              <tbody>
                {subdomains.map((subdomain, idx) => (
                  <tr key={subdomain.id || idx}>
                    <td className="font-medium">{subdomain.name}</td>
                    <td><span className="badge badge-info">{subdomain.source}</span></td>
                    <td><span className={`badge ${subdomain.is_active ? "badge-success" : "badge-secondary"}`}>{subdomain.is_active ? "Yes" : "No"}</span></td>
                    <td className="text-sm text-muted">{subdomain.created_at ? new Date(subdomain.created_at).toLocaleDateString() : "N/A"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {showFilters && renderPagination(subdomainPagination, subdomainPage, setSubdomainPage)}
        </>
      )}
    </>
  );

  const renderWHOIS = (whois) => {
    if (!whois) return <div className="p-4"><p className="text-muted">No WHOIS data found</p></div>;
    return (
      <div className="p-4 space-y-4">
        <div className="grid grid-2">
          <div><label className="text-sm font-semibold text-muted">Registrar</label><p>{whois.registrar || "N/A"}</p></div>
          <div><label className="text-sm font-semibold text-muted">Status</label><p>{whois.status || "N/A"}</p></div>
          <div><label className="text-sm font-semibold text-muted">Created Date</label><p>{whois.created_date ? new Date(whois.created_date).toLocaleDateString() : "N/A"}</p></div>
          <div><label className="text-sm font-semibold text-muted">Expiry Date</label><p>{whois.expiry_date ? new Date(whois.expiry_date).toLocaleDateString() : "N/A"}</p></div>
        </div>
        {whois.name_servers && (
          <div>
            <label className="text-sm font-semibold text-muted">Name Servers</label>
            <div className="flex flex-wrap gap-2 mt-2">
              {(() => {
                try {
                  const servers = typeof whois.name_servers === "string" ? JSON.parse(whois.name_servers) : whois.name_servers;
                  return Array.isArray(servers) ? servers.map((ns, idx) => <span key={idx} className="badge badge-info">{ns}</span>) : null;
                } catch (e) { return null; }
              })()}
            </div>
          </div>
        )}
        {whois.raw_data && (
          <div>
            <label className="text-sm font-semibold text-muted">Raw WHOIS Data</label>
            <pre className="mt-2 p-4 bg-gray-50 rounded-md text-xs overflow-x-auto">{whois.raw_data}</pre>
          </div>
        )}
      </div>
    );
  };

  // ============================================
  // SSL / PORT / IP RENDER
  // ============================================

  const renderJobCard = (job, renderDetail) => (
    <div key={job.id} className="p-4 border-b">
      <div className="flex items-center justify-between mb-2">
        <span className="text-sm text-muted">
          Scan ID: {job.id.slice(0, 8)}... | {new Date(job.created_at).toLocaleString()}
        </span>
        <button
          className={`btn btn-sm ${selectedJob?.id === job.id ? "btn-primary" : "btn-secondary"}`}
          onClick={() => handleViewJobResult(job)}
          disabled={loadingJobResults}
        >
          {selectedJob?.id === job.id ? "Hide" : "View Details"}
        </button>
      </div>
      {selectedJob?.id === job.id && (
        loadingJobResults ? (
          <div className="loading"><div className="spinner"></div><span>Loading...</span></div>
        ) : (
          renderDetail(selectedJobResults)
        )
      )}
    </div>
  );

  const renderSSLDetail = (results) => {
    if (!results?.results?.length) return <p className="text-muted p-2">No data</p>;
    return results.results.map((r, i) => (
      <div key={i} className="mt-2">
        <div className="flex items-center gap-2 mb-2">
          <span className={`badge ${r.grade === "A" ? "badge-success" : r.grade === "B" ? "badge-warning" : "badge-danger"}`}>
            Grade: {r.grade}
          </span>
          <strong>{r.domain}</strong>
        </div>
        <table className="table">
          <tbody>
            <tr><td><strong>Subject</strong></td><td>{r.certificate?.subject}</td></tr>
            <tr><td><strong>Issuer</strong></td><td>{r.certificate?.issuer}</td></tr>
            <tr><td><strong>Serial Number</strong></td><td>{r.certificate?.serial_number}</td></tr>
            <tr><td><strong>Valid From</strong></td><td>{r.certificate?.valid_from ? new Date(r.certificate.valid_from).toLocaleDateString() : "N/A"}</td></tr>
            <tr><td><strong>Valid Until</strong></td><td>{r.certificate?.valid_until ? new Date(r.certificate.valid_until).toLocaleDateString() : "N/A"}</td></tr>
            <tr><td><strong>Days Until Expiry</strong></td><td>{r.certificate?.days_until_expiry} days</td></tr>
            <tr><td><strong>Self Signed</strong></td><td>{r.certificate?.is_self_signed ? "Yes" : "No"}</td></tr>
            <tr><td><strong>TLS Version</strong></td><td>{r.connection?.tls_version}</td></tr>
            <tr><td><strong>Cipher Suite</strong></td><td>{r.connection?.cipher_suite}</td></tr>
            <tr><td><strong>SANs</strong></td><td>{r.certificate?.san?.join(", ")}</td></tr>
          </tbody>
        </table>
        {r.issues?.length > 0 && (
          <div className="alert alert-error mt-2">⚠️ Issues: {r.issues.join(", ")}</div>
        )}
      </div>
    ));
  };

  const renderPortDetail = (results) => {
    if (!results?.results?.length) return <p className="text-muted p-2">No data</p>;
    return results.results.map((r, i) => (
      <div key={i} className="mt-2">
        <p className="mb-2 text-sm">
          <strong>IP:</strong> {r.ip_address} &nbsp;|&nbsp;
          <strong>Scanned:</strong> {r.total_scanned} ports &nbsp;|&nbsp;
          <strong>Open:</strong> {r.open_ports?.length || 0} &nbsp;|&nbsp;
          <strong>Duration:</strong> {r.scan_duration_ms}ms
        </p>
        {!r.open_ports?.length ? (
          <p className="text-muted">No open ports found</p>
        ) : (
          <table className="table">
            <thead><tr><th>Port</th><th>Protocol</th><th>State</th><th>Service</th></tr></thead>
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

  const renderIPDetail = (results) => {
    if (!results?.results?.length) return <p className="text-muted p-2">No data</p>;
    return results.results.map((r, i) => (
      <div key={i} className="mt-2">
        <table className="table">
          <tbody>
            <tr><td><strong>IP Address</strong></td><td>{r.ip_address}</td></tr>
            <tr><td><strong>Country</strong></td><td>{r.geolocation?.country} ({r.geolocation?.country_code})</td></tr>
            <tr><td><strong>City / Region</strong></td><td>{r.geolocation?.city}, {r.geolocation?.region}</td></tr>
            <tr><td><strong>ISP</strong></td><td>{r.geolocation?.isp}</td></tr>
            <tr><td><strong>Organization</strong></td><td>{r.geolocation?.org}</td></tr>
            <tr><td><strong>ASN</strong></td><td>{r.asn?.number} — {r.asn?.name}</td></tr>
            <tr><td><strong>Reverse DNS</strong></td><td>{r.reverse_dns || "N/A"}</td></tr>
            <tr><td><strong>Coordinates</strong></td><td>{r.geolocation?.latitude}, {r.geolocation?.longitude}</td></tr>
          </tbody>
        </table>
      </div>
    ));
  };

  const renderSSLResults = (jobs) => {
    if (!jobs?.length) return <div className="p-4"><p className="text-muted">No SSL scan results found. Run an SSL scan first.</p></div>;
    return jobs.map((job) => renderJobCard(job, renderSSLDetail));
  };

  const renderPortResults = (jobs) => {
    if (!jobs?.length) return <div className="p-4"><p className="text-muted">No port scan results found. Run a port scan first.</p></div>;
    return jobs.map((job) => renderJobCard(job, renderPortDetail));
  };

  const renderIPResults = (jobs) => {
    if (!jobs?.length) return <div className="p-4"><p className="text-muted">No IP scan results found. Run an IP scan first.</p></div>;
    return jobs.map((job) => renderJobCard(job, renderIPDetail));
  };

  const renderAllResults = () => {
    const sslJobs = scanJobs.filter((j) => j.scan_type === "ssl");
    const portJobs = scanJobs.filter((j) => j.scan_type === "port");
    const ipJobs = scanJobs.filter((j) => j.scan_type === "ip");

    return (
      <div className="space-y-6">
        {/* WHOIS */}
        <div className="card">
          <div className="card-header">
            <h3 className="card-title flex items-center"><Server size={20} className="mr-2" />WHOIS Information</h3>
          </div>
          {renderWHOIS(whoisData)}
        </div>

        {/* DNS */}
        <div className="card">
          <div className="card-header">
            <h3 className="card-title flex items-center"><Globe size={20} className="mr-2" />DNS Records ({dnsPagination.total || 0})</h3>
          </div>
          {renderDNSRecords(dnsData, true)}
        </div>

        {/* Subdomains */}
        <div className="card">
          <div className="card-header">
            <h3 className="card-title flex items-center"><Globe size={20} className="mr-2" />Subdomains ({subdomainPagination.total || 0})</h3>
          </div>
          {renderSubdomains(subdomainData, true)}
        </div>

        {/* SSL */}
        {sslJobs.length > 0 && (
          <div className="card">
            <div className="card-header">
              <h3 className="card-title">🔒 SSL/TLS Results ({sslJobs.length})</h3>
            </div>
            {renderSSLResults(sslJobs)}
          </div>
        )}

        {/* Port */}
        {portJobs.length > 0 && (
          <div className="card">
            <div className="card-header">
              <h3 className="card-title">🔍 Port Scan Results ({portJobs.length})</h3>
            </div>
            {renderPortResults(portJobs)}
          </div>
        )}

        {/* IP */}
        {ipJobs.length > 0 && (
          <div className="card">
            <div className="card-header">
              <h3 className="card-title">🌍 IP Geolocation Results ({ipJobs.length})</h3>
            </div>
            {renderIPResults(ipJobs)}
          </div>
        )}
      </div>
    );
  };

  return (
    <div>
      <div className="page-header">
        <h1 className="page-title">Scan Results</h1>
        <p className="page-description">View and analyze reconnaissance data collected from scans</p>
      </div>

      {error && <div className="alert alert-error mb-4">{error}</div>}

      {/* Filters */}
      <div className="card mb-4">
        <div className="grid grid-2 gap-4">
          <div className="form-group">
            <label className="form-label">Select Asset</label>
            <select className="form-select" value={selectedAsset} onChange={(e) => setSelectedAsset(e.target.value)}>
              {assets.length === 0 ? (
                <option>No assets available</option>
              ) : (
                assets.map((asset) => (
                  <option key={asset.id} value={asset.id}>{asset.name} ({asset.type})</option>
                ))
              )}
            </select>
          </div>

          <div className="form-group">
            <label className="form-label">Result Type</label>
            <select className="form-select" value={resultType} onChange={(e) => setResultType(e.target.value)}>
              <option value="all">All Results</option>
              <option value="dns">DNS Records</option>
              <option value="subdomains">Subdomains</option>
              <option value="whois">WHOIS Information</option>
              <option value="ssl">🔒 SSL/TLS Scans</option>
              <option value="port">🔍 Port Scans</option>
              <option value="ip">🌍 IP Geolocation</option>
            </select>
          </div>
        </div>

        {selectedAssetData && (
          <div className="mt-4 p-4 bg-gray-50 rounded-md">
            <h4 className="font-semibold text-sm mb-2">Asset Details:</h4>
            <div className="flex items-center gap-4 text-sm text-muted">
              <span><strong>Name:</strong> {selectedAssetData.name}</span>
              <span><strong>Type:</strong> {selectedAssetData.type}</span>
              <span><strong>Status:</strong> {selectedAssetData.status}</span>
            </div>
          </div>
        )}
      </div>

      {/* Results Display */}
      {loading ? (
        <div className="card">
          <div className="loading"><div className="spinner"></div><span>Loading results...</span></div>
        </div>
      ) : !selectedAsset ? (
        <div className="card">
          <div className="empty-state">
            <FileText className="empty-state-icon" size={64} />
            <h3 className="empty-state-title">No asset selected</h3>
            <p className="empty-state-description">Select an asset to view its scan results</p>
          </div>
        </div>
      ) : (
        <>
          {resultType === "all" && renderAllResults()}

          {resultType === "dns" && (
            <div className="card">
              <div className="card-header"><h3 className="card-title"><Globe size={20} className="inline mr-2" />DNS Records</h3></div>
              {renderDNSRecords(dnsData, true)}
            </div>
          )}

          {resultType === "subdomains" && (
            <div className="card">
              <div className="card-header"><h3 className="card-title"><Globe size={20} className="inline mr-2" />Subdomains</h3></div>
              {renderSubdomains(subdomainData, true)}
            </div>
          )}

          {resultType === "whois" && (
            <div className="card">
              <div className="card-header"><h3 className="card-title"><Server size={20} className="inline mr-2" />WHOIS Information</h3></div>
              {renderWHOIS(whoisData)}
            </div>
          )}

          {resultType === "ssl" && (
            <div className="card">
              <div className="card-header"><h3 className="card-title">🔒 SSL/TLS Scan Results</h3></div>
              {renderSSLResults(scanJobs.filter((j) => j.scan_type === "ssl"))}
            </div>
          )}

          {resultType === "port" && (
            <div className="card">
              <div className="card-header"><h3 className="card-title">🔍 Port Scan Results</h3></div>
              {renderPortResults(scanJobs.filter((j) => j.scan_type === "port"))}
            </div>
          )}

          {resultType === "ip" && (
            <div className="card">
              <div className="card-header"><h3 className="card-title">🌍 IP Geolocation Results</h3></div>
              {renderIPResults(scanJobs.filter((j) => j.scan_type === "ip"))}
            </div>
          )}
        </>
      )}
    </div>
  );
}

export default Results;
