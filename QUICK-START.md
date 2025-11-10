# 🚀 QUICK START GUIDE - Frontend

## Prerequisites
- Node.js 18+ installed
- Backend running on http://localhost:8080

## Installation Steps

### 1. Install Dependencies
```bash
cd yayasan-erp/frontend
npm install
```

This will install:
- React 18
- TypeScript
- Vite
- Tailwind CSS
- React Router
- React Query
- Zustand
- Axios
- Lucide Icons
- And more...

### 2. Configure Environment
The `.env` file is already configured with defaults:
```env
VITE_API_URL=http://localhost:8080/api/v1
```

If your backend runs on a different port, edit `.env` file.

### 3. Start Development Server
```bash
npm run dev
```

The app will start on **http://localhost:3000**

### 4. Login
```
Username: admin
Password: admin123
```

(Make sure you've run the backend seed script first!)

## Available Scripts

```bash
# Development
npm run dev              # Start dev server (hot reload)

# Build
npm run build           # Build for production
npm run preview         # Preview production build

# Linting
npm run lint            # Run ESLint
```

## Project Structure

```
src/
├── components/
│   ├── common/         # Reusable UI components
│   │   ├── Button.tsx
│   │   ├── Input.tsx
│   │   ├── Card.tsx
│   │   ├── Table.tsx
│   │   ├── Badge.tsx
│   │   └── Modal.tsx
│   └── layout/         # Layout components
│
├── pages/              # Page components
│   ├── auth/
│   │   └── LoginPage.tsx        ✅ Complete
│   ├── dashboard/
│   │   └── DashboardPage.tsx    ✅ Complete
│   └── students/
│       └── StudentsPage.tsx     ✅ Complete
│
├── layouts/
│   └── MainLayout.tsx           ✅ Complete (with sidebar)
│
├── services/           # API services
│   ├── api.ts                   ✅ Axios instance
│   └── auth.service.ts          ✅ Auth API
│
├── store/              # State management
│   └── auth.store.ts            ✅ Zustand store
│
├── types/              # TypeScript types
│   └── index.ts                 ✅ All types defined
│
├── utils/              # Utilities
│   ├── cn.ts                    ✅ Class merger
│   └── format.ts                ✅ Formatters
│
├── App.tsx             # Main app with routing
├── main.tsx            # Entry point
└── index.css           # Global styles

## Features Implemented

### ✅ Authentication
- Login page with form validation
- JWT token management
- Auto token injection to API calls
- Protected routes
- Logout functionality

### ✅ Dashboard
- Statistics cards (students, revenue, expenses, etc)
- Recent activities
- Quick actions
- Responsive layout

### ✅ Students Management
- Student list with pagination
- Search functionality
- Data table with sorting
- Status badges
- Responsive design

### ✅ Layout & Navigation
- Sidebar navigation with submenu
- Responsive mobile menu
- User profile display
- Logout button
- Branch selector

### ✅ UI Components
- Button (multiple variants)
- Input (with validation)
- Card
- Table
- Badge
- Modal
- Loading states

### ✅ Utilities
- Currency formatter (IDR)
- Date formatter (Indonesian)
- Number formatter
- Class name merger (Tailwind)

## Next Steps - Build More Pages

### 1. Student Detail Page
Create `src/pages/students/StudentDetailPage.tsx`:
- View student information
- Parent/guardian details
- Academic history
- Invoices & payments

### 2. Invoice Page
Create `src/pages/students/InvoicesPage.tsx`:
- List all invoices
- Create new invoice
- Invoice detail
- Print invoice

### 3. Payment Page
Create `src/pages/students/PaymentsPage.tsx`:
- Record payments
- Payment history
- Receipt printing

### 4. Finance Pages
Create finance module pages:
- Chart of Accounts (tree view)
- Journal entries (with approval)
- Financial reports (viewer)
- Budget management

### 5. Employee Pages
Create HR module pages:
- Employee list
- Employee detail
- Payroll processing
- Attendance

### 6. Asset Pages
Create asset module pages:
- Asset list
- Asset detail
- Maintenance schedule
- Depreciation report

### 7. Inventory Pages
Create inventory pages:
- Item list
- Stock transactions
- Stock opname
- Reports

## Tips for Development

### 1. Using React Query
```typescript
const { data, isLoading, error } = useQuery({
  queryKey: ['students', page],
  queryFn: async () => {
    const response = await api.get('/students', {
      params: { page }
    });
    return response.data.data;
  },
});
```

### 2. Using Zustand Store
```typescript
const { user, setUser } = useAuthStore();
```

### 3. Formatting
```typescript
import { formatCurrency, formatDate } from '@/utils/format';

formatCurrency(1500000)  // Rp 1.500.000
formatDate('2025-01-15') // 15 Januari 2025
```

### 4. Using Components
```typescript
import { Button, Card, Table } from '@/components/common';

<Button variant="primary" onClick={handleClick}>
  Save
</Button>

<Card title="Students">
  <Table columns={columns} data={students} />
</Card>
```

## Common Issues

### Port 3000 already in use
```bash
# Kill process on port 3000
# Linux/Mac:
lsof -ti:3000 | xargs kill -9

# Windows:
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Or use different port:
npm run dev -- --port 3001
```

### Module not found
```bash
npm install
npm run dev
```

### API not connecting
1. Check backend is running: http://localhost:8080/health
2. Check VITE_API_URL in .env
3. Check CORS configuration in backend

### Build errors
```bash
# Clear cache
rm -rf node_modules package-lock.json
npm install
npm run build
```

## Deployment

### Build for Production
```bash
npm run build
```

Output will be in `dist/` folder.

### Deploy to Vercel
```bash
npm install -g vercel
vercel
```

### Deploy to Netlify
```bash
npm install -g netlify-cli
netlify deploy
```

### Deploy to Nginx
```bash
npm run build

# Copy dist/ to server
scp -r dist/* user@server:/var/www/yayasan-erp/

# Nginx config:
server {
    listen 80;
    server_name yourdomain.com;
    root /var/www/yayasan-erp;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
    }
}
```

## 🎉 You're Ready!

Start building amazing features! 🚀

For questions, check:
- README.md (main documentation)
- Component source code
- React Query docs: https://tanstack.com/query
- Tailwind docs: https://tailwindcss.com
