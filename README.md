# Yayasan ERP - Frontend Web Application

Modern React + TypeScript frontend for Yayasan ERP System.

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ and npm
- Backend API running on http://localhost:8080

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

The app will run on **http://localhost:3000**

## 📁 Project Structure

```
src/
├── components/          # Reusable UI components
│   ├── common/         # Buttons, inputs, cards, etc
│   ├── layout/         # Sidebar, header, footer
│   ├── auth/           # Login form, etc
│   └── dashboard/      # Dashboard widgets
├── pages/              # Page components
│   ├── auth/           # Login, register pages
│   ├── dashboard/      # Dashboard page
│   ├── students/       # Student management pages
│   ├── finance/        # Finance pages
│   ├── employees/      # HR pages
│   ├── assets/         # Asset management pages
│   └── inventory/      # Inventory pages
├── services/           # API service layers
├── hooks/              # Custom React hooks
├── types/              # TypeScript type definitions
├── utils/              # Utility functions
├── store/              # Zustand state management
└── layouts/            # Layout components
```

## 🎨 Tech Stack

- **Framework:** React 18 + TypeScript
- **Build Tool:** Vite
- **Styling:** Tailwind CSS
- **State Management:** Zustand + React Query
- **Routing:** React Router v6
- **Forms:** React Hook Form + Zod
- **HTTP Client:** Axios
- **Icons:** Lucide React
- **Charts:** Recharts

## 🔑 Default Login

```
Username: admin
Password: admin123
```

## 📋 Available Pages (To Be Created)

### Core Pages
- [x] Login Page
- [ ] Dashboard
- [ ] Profile

### Student Management
- [ ] Student List
- [ ] Student Detail
- [ ] Add/Edit Student
- [ ] Student Invoices
- [ ] Payment Processing

### Finance & Accounting
- [ ] Chart of Accounts
- [ ] Journal Entries
- [ ] Financial Reports
- [ ] Budget Management

### HR & Payroll
- [ ] Employee List
- [ ] Employee Detail
- [ ] Payroll Processing
- [ ] Attendance
- [ ] Leave Management

### Asset Management
- [ ] Asset List
- [ ] Asset Detail
- [ ] Maintenance Schedule
- [ ] Depreciation Report

### Inventory
- [ ] Inventory Items
- [ ] Stock Transactions
- [ ] Stock Opname
- [ ] Purchase Orders

## 🎯 Features

### Implemented
- ✅ API Service with interceptors
- ✅ Authentication (Login/Logout)
- ✅ Protected Routes
- ✅ Token Management
- ✅ Type Definitions
- ✅ Utility Functions
- ✅ State Management

### To Be Implemented
- [ ] Dashboard with statistics
- [ ] CRUD operations for all modules
- [ ] Form validation
- [ ] Data tables with pagination
- [ ] Charts and graphs
- [ ] Export to Excel/PDF
- [ ] Print functionality
- [ ] Dark mode
- [ ] Multi-language

## 🔧 Environment Variables

Create `.env` file:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 📝 API Integration

All API calls go through `/src/services/api.ts` which:
- Adds Bearer token to requests
- Handles 401 (Unauthorized) redirects
- Provides error handling
- Manages request/response interceptors

## 🎨 Styling

Using **Tailwind CSS** with custom configuration:
- Primary color palette
- Responsive design
- Dark mode ready
- Custom components

## 🔐 Authentication Flow

1. User submits login form
2. Frontend calls `/api/v1/auth/login`
3. Backend returns JWT token + user data
4. Token stored in localStorage
5. Token added to all subsequent requests
6. On 401 error, redirect to login

## 📦 Build & Deploy

```bash
# Build for production
npm run build

# Output directory: dist/
# Deploy to: Nginx, Apache, Vercel, Netlify, etc.
```

### Nginx Configuration Example

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    root /var/www/yayasan-erp/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 🐛 Troubleshooting

### CORS Issues
Make sure backend allows frontend origin:
```go
// In backend: internal/routes/routes.go
router.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:3000"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Authorization", "Content-Type"},
}))
```

### API Connection Failed
1. Check backend is running on port 8080
2. Check VITE_API_URL in .env
3. Check browser console for errors

## 📚 Next Steps

1. **Create Components:**
   - Button, Input, Card, Table
   - Modal, Alert, Loading

2. **Create Pages:**
   - Dashboard with stats
   - Student list with filters
   - Invoice generation
   - Payroll processing

3. **Add Features:**
   - Data export
   - Print functionality
   - Advanced filters
   - Bulk operations

## 🤝 Contributing

1. Create feature branch
2. Make changes
3. Test thoroughly
4. Submit PR

## 📄 License

Private - Yayasan ERP System

---

**Need help?** Check the API documentation or contact support.
