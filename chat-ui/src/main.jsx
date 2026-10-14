import { createRoot } from 'react-dom/client'
import App from './App.jsx'
import './styles.css'

// No StrictMode on purpose: the backend allows one socket per user,
// and StrictMode's double-mount would open and drop an extra one.
createRoot(document.getElementById('root')).render(<App />)
