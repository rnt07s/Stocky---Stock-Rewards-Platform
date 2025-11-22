# Git & GitHub Setup Guide

## Initialize Git Repository

```powershell
cd "c:\Users\RAUNEET SINGH\OneDrive\Desktop\backend Stocky"

# Initialize Git repository
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: Stocky stock rewards platform

- Implemented REST API with Gin framework
- PostgreSQL database with double-entry ledger
- Stock price service with hourly updates
- Idempotency and edge case handling
- Complete API documentation and setup guides
"
```

## Create GitHub Repository

### Option 1: Using GitHub Web Interface

1. Go to https://github.com/new
2. Repository name: `stocky-assignment`
3. Description: `Stock Rewards Platform - Award fractional shares of Indian stocks as incentives`
4. Keep it **Public** (for assignment submission)
5. **Do NOT** initialize with README (we already have one)
6. Click "Create repository"

### Option 2: Using GitHub CLI (if installed)

```powershell
# Install GitHub CLI from: https://cli.github.com/
gh auth login
gh repo create stocky-assignment --public --source=. --remote=origin --push
```

## Push to GitHub

After creating the repository on GitHub, you'll see instructions. Run these commands:

```powershell
# Add remote repository (replace YOUR_USERNAME with your GitHub username)
git remote add origin https://github.com/YOUR_USERNAME/stocky-assignment.git

# Rename branch to main (if needed)
git branch -M main

# Push code to GitHub
git push -u origin main
```

## Verify Upload

Visit your repository URL:
```
https://github.com/YOUR_USERNAME/stocky-assignment
```

You should see:
- ✅ All source code files
- ✅ README.md displayed on homepage
- ✅ Complete directory structure
- ✅ Green checkmark (Git commit successful)

## Repository Structure on GitHub

```
stocky-assignment/
├── .gitignore
├── README.md                              # Main documentation
├── ARCHITECTURE.md                        # Design decisions
├── WINDOWS_SETUP.md                       # Windows setup guide
├── GITHUB_SETUP.md                        # This file
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── setup.ps1
├── test-api.ps1
├── test-api.sh
├── Stocky-API.postman_collection.json
├── .env.example
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   └── services/
└── migrations/
    └── 001_initial_schema.sql
```

## Best Practices for Commit Messages

Follow conventional commits:

```powershell
# Feature
git commit -m "feat: add user portfolio endpoint"

# Bug fix
git commit -m "fix: resolve rounding error in fee calculation"

# Documentation
git commit -m "docs: update API examples in README"

# Refactor
git commit -m "refactor: extract fee calculation to separate function"

# Test
git commit -m "test: add unit tests for reward service"
```

## Making Changes After Initial Push

```powershell
# Make your changes to files
# ...

# Stage changes
git add .

# Commit with descriptive message
git commit -m "fix: improve error handling in stock price service"

# Push to GitHub
git push
```

## Adding a Professional Touch

### 1. Add GitHub Repository Topics

On GitHub repository page:
1. Click "Settings" → "General"
2. In "Topics" section, add:
   - `golang`
   - `gin-framework`
   - `postgresql`
   - `rest-api`
   - `stock-trading`
   - `fintech`
   - `double-entry-bookkeeping`

### 2. Add Repository Description

On GitHub repository page:
1. Click "⚙️ Settings" → "General"
2. Set description:
   ```
   🏦 Stock Rewards Platform - Award fractional shares of Indian stocks (NSE/BSE) as user incentives. Built with Go, Gin, PostgreSQL. Features: double-entry ledger, idempotency, hourly price updates, comprehensive API.
   ```

### 3. Add LICENSE File (Optional)

```powershell
# Create MIT License file
@"
MIT License

Copyright (c) 2025 [Your Name]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
"@ | Out-File -FilePath LICENSE -Encoding utf8

git add LICENSE
git commit -m "docs: add MIT license"
git push
```

### 4. Pin Repository

On your GitHub profile:
1. Go to your profile page
2. Click "Customize your pins"
3. Select `stocky-assignment`
4. This makes it visible to recruiters/evaluators

## Troubleshooting

### Error: "failed to push some refs"

```powershell
# Pull latest changes first
git pull origin main --rebase

# Then push
git push
```

### Error: "remote origin already exists"

```powershell
# Remove existing remote
git remote remove origin

# Add correct remote
git remote add origin https://github.com/YOUR_USERNAME/stocky-assignment.git
```

### Large Files Warning

If you get warnings about large files:
```powershell
# Remove from tracking (if accidentally added)
git rm --cached go.sum
git commit -m "chore: remove go.sum from tracking"

# Then re-add properly
git add go.sum
git commit -m "chore: add go.sum"
```

## Assignment Submission

When submitting the assignment, provide:

**GitHub Repository URL**:
```
https://github.com/YOUR_USERNAME/stocky-assignment
```

**Key Files to Highlight**:
1. `README.md` - Complete API documentation
2. `ARCHITECTURE.md` - Design decisions and edge cases
3. `migrations/001_initial_schema.sql` - Database schema
4. `cmd/server/main.go` - Application entry point
5. `internal/services/services.go` - Business logic

**Live Demo** (Optional):
If deployed to cloud (Heroku, Railway, Render):
```
Live API: https://stocky-api.herokuapp.com
Health Check: https://stocky-api.herokuapp.com/health
```

## Final Checklist

Before submission, verify:

- ✅ All code is pushed to GitHub
- ✅ README.md renders correctly on GitHub
- ✅ .env file is NOT pushed (check .gitignore)
- ✅ Repository is Public
- ✅ No sensitive credentials in code
- ✅ API endpoints documented with examples
- ✅ Database schema included
- ✅ Edge cases explained
- ✅ Setup instructions clear and complete

---

**Good luck with your assignment! 🚀**
