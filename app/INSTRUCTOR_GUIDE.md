# 🎯 QUICK START GUIDE - Instructor

## Chuẩn Bị Trước Khóa Học

### 1. Clone/Setup Repository

```bash
cd /path/to/CMC-InternProgram/app
```

### 2. Verify Structure

```
week1/
├── session1-foundation/      ✅ Done
├── session2-basic-api/       ✅ Done (key files)
├── session3-database/        ✅ README done
├── session4-advanced/        ✅ README done
├── session5-testing/         ✅ README done
├── session6-deployment/      ✅ README done
├── docs/                     ✅ Existing
│   ├── api.yml
│   ├── requirements.md
│   └── scope.md
└── TRAINING_PLAN.md          ✅ Done
```

### 3. Cần Hoàn Thiện

Session 2-6 đã có:

- ✅ README chi tiết với teaching notes
- ✅ Code snippets quan trọng
- ✅ Teaching flow rõ ràng

Cần thêm (optional - có thể code live hoặc dùng snippets):

- Full implementation files (có thể copy từ session2 template)
- go.mod, go.sum cho mỗi session
- Migration files cho session 3
- Test files cho session 5

## Teaching Strategy

### Approach: "Code Walkthrough" (Không Code Từ Đầu)

Mỗi session:

1. **Show code có sẵn** từ folder sessionX
2. **Explain từng phần** với teaching notes trong code
3. **So sánh với session trước** - highlight changes
4. **Demo chạy** - show kết quả thực tế
5. **Q&A** - interactive discussion

### Time Allocation per Session (3h)

```
0:00 - 0:15   Review session trước + Q&A
0:15 - 0:30   Theory/Concepts (slides)
0:30 - 2:30   Code Walkthrough (chính)
2:30 - 2:45   Live Demo
2:45 - 3:00   Homework explanation + wrap-up
```

## Buổi 1: Foundation ⚡

**Mindset:** Simple start, foundation cho tương lai

**Key Messages:**

- Project structure matters
- Clean code from day 1
- Git workflow

**Demo:**

```bash
cd session1-foundation
go run cmd/server/main.go
curl http://localhost:8080/health
```

**Homework:** Add `/hello/{name}` endpoint

---

## Buổi 2: Basic API 🏗️

**Mindset:** Clean Architecture in action

**Key Messages:**

- 4 layers: Model → Storage → Service → Handler
- Dependency injection
- Interface for flexibility

**Comparison Point:**

```
Buổi 1: 1 file, everything mixed
Buổi 2: 4 layers, separation of concerns
```

**Demo Flow:**

1. Show folder structure
2. Bottom-up: Model → Storage → Service → Handler
3. Wire up in main.go
4. Test API with curl/Postman

**Critical Teaching Moments:**

- Why interface? → Buổi 3 sẽ swap implementation!
- Why layers? → Easy to test, maintain, scale

---

## Buổi 3: Database 🗄️

**Mindset:** Persistence without pain

**Key Messages:**

- Clean Architecture pays off: CHỈ 1-2 dòng thay đổi!
- Docker for local development
- Migrations for schema management

**Demo Flow:**

1. Review memory storage from Buổi 2
2. Restart server → data lost!
3. Show PostgreSQL implementation
4. Compare interfaces (same!)
5. Change main.go: 1 line!
6. Test persistence

**"Aha Moment":**

```go
// Buổi 2
store := memory.NewMemoryStorage()

// Buổi 3 - CHỈ ĐỔI 1 DÒNG!
store := postgres.NewPostgresStorage(db)

// Handler, Service, Model: KHÔNG ĐỔI!
```

---

## Buổi 4: Advanced Features 🔍

**Mindset:** Production-ready features

**Key Messages:**

- Pagination for scalability
- Validation for security
- Sorting for UX

**Demo:** Complex query

```bash
curl "localhost:8080/assets?type=domain&status=active&search=example&page=1&sort_by=name"
```

---

## Buổi 5: Testing 🧪

**Mindset:** "If it's not tested, it's broken"

**Key Messages:**

- Unit tests with mocks
- Integration tests for confidence
- Table-driven tests pattern
- Coverage as metric

**Demo:**

```bash
go test ./... -cover
go tool cover -html=coverage.out
```

**Show:** Green tests = confidence to refactor

---

## Buổi 6: Deployment 🚢

**Mindset:** From localhost to production

**Key Messages:**

- Frontend integration (CORS!)
- Docker for consistency
- Docker Compose for orchestration

**Demo:** Full stack

```bash
docker-compose up -d
# Open http://localhost
# Show frontend working
# Show database persisting
# Stop/restart → still works!
```

**"Wow Moment":** Entire stack in 1 command!

---

## Teaching Tips

### Do's ✅

- Use visual diagrams on board
- Ask questions frequently: "Tại sao...?"
- Show mistakes and how to fix
- Compare with bad practices
- Celebrate "aha moments"
- Use real-world analogies

### Don'ts ❌

- Don't code everything from scratch (too slow)
- Don't skip theory (students need concepts)
- Don't rush through errors (learning moments!)
- Don't assume knowledge (explain basics)

### Engagement Techniques

- Pair programming for homework review
- Group discussions for architecture decisions
- Live debugging sessions
- Code review practice

---

## Common Student Questions & Answers

**Q: Tại sao không dùng framework (gin, echo)?**
A: Learn fundamentals first! Framework hides complexity. Sau khi hiểu stdlib, framework dễ dàng.

**Q: Clean Architecture có quá phức tạp không?**
A: Yes for toy projects. Perfect for real projects. Show evolution: Buổi 1 → Buổi 6.

**Q: Khi nào dùng MVC vs Clean Architecture?**
A: Show comparison table trong CLEAN_ARCHITECTURE.MD section 5.

**Q: Có cần test 100% coverage không?**
A: No. Focus on critical paths. 70-80% is good.

**Q: Production deployment khác gì local?**
A: Show in Buổi 6: HTTPS, secrets management, monitoring, scaling.

---

## Troubleshooting Guide

### Issue: Port already in use

```bash
lsof -i :8080
kill -9 <PID>
```

### Issue: Docker not starting

```bash
docker ps
docker-compose down
docker system prune
```

### Issue: Go module errors

```bash
go mod tidy
go clean -modcache
```

### Issue: Database connection failed

- Check docker-compose up
- Check .env file
- Check DB_HOST environment variable

---

## Resources for Students

**Must Read:**

- [ ] CLEAN_ARCHITECTURE.MD (full document)
- [ ] Each session's README before class

**Practice:**

- [ ] Complete homework before next session
- [ ] Experiment with code
- [ ] Break things and fix them!

**Extra:**

- Go by Example
- Effective Go
- Clean Code (Robert Martin)

---

## Success Metrics

After training, students should:

- ✅ Build CRUD API independently
- ✅ Explain Clean Architecture benefits
- ✅ Write basic tests
- ✅ Deploy with Docker
- ✅ Debug and troubleshoot
- ✅ Read and understand others' code

---

## Feedback & Improvement

After each session:

1. Collect feedback (5 min survey)
2. Note difficult topics
3. Adjust next session pace
4. Update teaching notes

---

## 🎓 Final Notes for Instructor

**Remember:**

- Students learn by DOING, not just watching
- Mistakes are learning opportunities
- Every question is valuable
- Patience and encouragement matter
- Real-world context helps retention

**Your Goal:**
Not just teach Go and Clean Architecture, but **instill good software engineering practices** that last a career.

**Good luck! 🚀**

---

## Quick Command Reference

```bash
# Session 1
cd session1-foundation && go run cmd/server/main.go

# Session 2
cd session2-basic-api && go run cmd/server/main.go

# Session 3
cd session3-database
docker-compose up -d
go run cmd/server/main.go

# Session 4
# Same as session 3 with more features

# Session 5
cd session5-testing
go test ./... -cover

# Session 6
cd session6-deployment
docker-compose up --build
# Open http://localhost
```

---

**Prepared by:** CMC Intern Program  
**Last Updated:** March 3, 2026  
**Version:** 1.0
