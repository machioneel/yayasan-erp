# ✅ TESTING CHECKLIST - Frontend

## Pre-Testing Setup

### ✅ Backend Running
- [ ] PostgreSQL running
- [ ] Backend server running on :8080
- [ ] Database seeded with initial data
- [ ] Health check: http://localhost:8080/health returns OK

### ✅ Frontend Setup
- [ ] Dependencies installed (`npm install`)
- [ ] .env file configured
- [ ] Dev server running on :3000

---

## 1. LOGIN PAGE TESTING

### Test 1.1: Page Load
- [ ] Navigate to http://localhost:3000
- [ ] Page redirects to /login
- [ ] Login form displays correctly
- [ ] Logo and title visible
- [ ] Demo credentials box showing

### Test 1.2: Form Validation
- [ ] Try submit with empty fields → Should show validation
- [ ] Enter wrong credentials → Should show error message
- [ ] Error message displays properly in red box

### Test 1.3: Successful Login
- [ ] Enter: admin / admin123
- [ ] Click "Masuk" button
- [ ] Button shows loading spinner
- [ ] Redirects to /dashboard on success
- [ ] Token saved in localStorage

**Expected Result:** ✅ Login successful, redirected to dashboard

---

## 2. DASHBOARD TESTING

### Test 2.1: Dashboard Display
- [ ] Dashboard page loads
- [ ] Statistics cards display (8 cards)
- [ ] All numbers show correctly
- [ ] Icons display for each card
- [ ] Recent activities section shows
- [ ] Quick actions section shows

### Test 2.2: Statistics Cards
- [ ] Total Siswa card shows number
- [ ] Pendapatan card shows currency (Rp format)
- [ ] Pengeluaran card shows currency
- [ ] Laba Bersih card shows currency
- [ ] Invoice Pending shows count
- [ ] Invoice Overdue shows count (in red)
- [ ] Total Karyawan shows count
- [ ] Stock Menipis shows count (in orange)

### Test 2.3: Navigation
- [ ] Sidebar visible on left
- [ ] User info displayed at top
- [ ] Menu items clickable
- [ ] Logout button at bottom

**Expected Result:** ✅ Dashboard fully functional

---

## 3. NAVIGATION TESTING

### Test 3.1: Sidebar Menu
- [ ] Click Dashboard → navigates to /dashboard
- [ ] Click Siswa → expands submenu
  - [ ] Daftar Siswa
  - [ ] Pendaftaran Baru
  - [ ] Tagihan & Invoice
  - [ ] Pembayaran
- [ ] Click Keuangan → expands submenu
- [ ] Click SDM & Payroll → expands submenu
- [ ] Click Aset → expands submenu
- [ ] Click Inventori → expands submenu

### Test 3.2: Active States
- [ ] Current page highlighted in blue
- [ ] Submenu items show when expanded
- [ ] Smooth transitions

### Test 3.3: Mobile Menu
- [ ] Resize to mobile width
- [ ] Sidebar hides
- [ ] Hamburger menu appears
- [ ] Click hamburger → sidebar slides in
- [ ] Click overlay → sidebar closes

**Expected Result:** ✅ Navigation works perfectly

---

## 4. STUDENTS PAGE TESTING

### Test 4.1: Page Load
- [ ] Click "Siswa" → "Daftar Siswa"
- [ ] Page loads at /students
- [ ] "Data Siswa" title shows
- [ ] "Tambah Siswa" button visible
- [ ] Search box visible
- [ ] Filter and Export buttons visible

### Test 4.2: Table Display
- [ ] Table shows student data
- [ ] Columns: No Registrasi, Nama, JK, Kelas, Cabang, Status, Tgl Daftar
- [ ] Status badges colored correctly:
  - Active → Green
  - Inactive → Gray
  - Graduated → Blue
  - Dropped → Red
- [ ] Date formatted correctly (Indonesian format)

### Test 4.3: Search Function
- [ ] Type in search box
- [ ] Results filter (when connected to backend)
- [ ] Clear search → shows all data

### Test 4.4: Pagination
- [ ] Pagination controls at bottom
- [ ] Shows "Menampilkan X - Y dari Z siswa"
- [ ] Previous button disabled on first page
- [ ] Next button disabled on last page
- [ ] Click Next → loads page 2
- [ ] Click Previous → back to page 1

### Test 4.5: Empty State
- [ ] When no data: shows "Tidak ada data siswa"
- [ ] Message centered

**Expected Result:** ✅ Students page functional

---

## 5. AUTHENTICATION TESTING

### Test 5.1: Protected Routes
- [ ] Logout from system
- [ ] Try to access /dashboard directly
- [ ] Should redirect to /login
- [ ] Try to access /students
- [ ] Should redirect to /login

### Test 5.2: Token Management
- [ ] Login successfully
- [ ] Check localStorage for 'token'
- [ ] Check localStorage for 'user'
- [ ] Token included in API requests (check Network tab)

### Test 5.3: Logout
- [ ] Click logout button
- [ ] Redirects to /login
- [ ] Token removed from localStorage
- [ ] User data removed from localStorage
- [ ] Cannot access protected routes

**Expected Result:** ✅ Auth working correctly

---

## 6. API INTEGRATION TESTING

### Test 6.1: API Calls
Open Browser DevTools → Network tab

#### Login API
- [ ] POST /api/v1/auth/login
- [ ] Status: 200
- [ ] Response contains: token, user data
- [ ] Token saved

#### Students API (when implemented)
- [ ] GET /api/v1/students?page=1
- [ ] Authorization header present
- [ ] Status: 200
- [ ] Response contains: data array, pagination

### Test 6.2: Error Handling
- [ ] Stop backend
- [ ] Try to access students page
- [ ] Shows loading state first
- [ ] Then shows error or empty state
- [ ] No crashes

### Test 6.3: 401 Handling
- [ ] Remove token from localStorage
- [ ] Make any API call
- [ ] Should redirect to login

**Expected Result:** ✅ API integration works

---

## 7. UI/UX TESTING

### Test 7.1: Responsiveness
Test on different screen sizes:
- [ ] Desktop (1920x1080) → All good
- [ ] Laptop (1366x768) → All good
- [ ] Tablet (768x1024) → Sidebar adapts
- [ ] Mobile (375x667) → Mobile menu works

### Test 7.2: Loading States
- [ ] Login shows spinner when submitting
- [ ] Dashboard shows spinner when loading
- [ ] Table shows spinner when loading
- [ ] Buttons disabled during loading

### Test 7.3: Colors & Theme
- [ ] Primary color: Blue (#3b82f6)
- [ ] Success: Green
- [ ] Warning: Yellow
- [ ] Danger: Red
- [ ] Text readable
- [ ] Contrast good

### Test 7.4: Icons
- [ ] All icons display (Lucide React)
- [ ] Icons sized correctly
- [ ] Icons colored correctly

**Expected Result:** ✅ UI/UX excellent

---

## 8. BROWSER COMPATIBILITY

### Test on Different Browsers
- [ ] Chrome/Edge (Chromium) → ✅ Primary
- [ ] Firefox → Should work
- [ ] Safari → Should work
- [ ] Mobile browsers → Should work

**Expected Result:** ✅ Works on modern browsers

---

## 9. PERFORMANCE TESTING

### Test 9.1: Load Time
- [ ] Initial load < 2 seconds
- [ ] Page transitions smooth
- [ ] No lag when typing

### Test 9.2: Bundle Size
```bash
npm run build
# Check dist/ size
```
- [ ] dist/ folder < 5 MB
- [ ] Main JS bundle < 500 KB (gzipped)

**Expected Result:** ✅ Fast performance

---

## 10. CONSOLE CHECK

### Test 10.1: No Errors
- [ ] Open Console (F12)
- [ ] Navigate through pages
- [ ] No red errors
- [ ] No warnings (or minimal)

### Test 10.2: Network
- [ ] All API calls successful
- [ ] No 404 errors
- [ ] Assets loading correctly

**Expected Result:** ✅ Clean console

---

## FINAL CHECKLIST

### Core Functionality
- [ ] ✅ Login works
- [ ] ✅ Logout works
- [ ] ✅ Protected routes work
- [ ] ✅ Dashboard displays
- [ ] ✅ Navigation works
- [ ] ✅ Students page works
- [ ] ✅ API calls work
- [ ] ✅ Responsive design
- [ ] ✅ No console errors

### UI Components
- [ ] ✅ Button component works
- [ ] ✅ Input component works
- [ ] ✅ Card component works
- [ ] ✅ Table component works
- [ ] ✅ Badge component works
- [ ] ✅ Modal component ready

### Ready for Development
- [ ] ✅ All dependencies installed
- [ ] ✅ TypeScript configured
- [ ] ✅ Tailwind configured
- [ ] ✅ React Query configured
- [ ] ✅ Routing configured
- [ ] ✅ State management ready

---

## 🎉 TESTING COMPLETE!

If all checkboxes are ✅, your frontend is ready for development!

### Next Steps:
1. Start building more pages
2. Connect to real backend APIs
3. Add more features
4. Deploy to production

**Status:** 🚀 READY FOR DEVELOPMENT!
