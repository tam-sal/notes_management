import './App.css'

import Header from './components/Header/Header'
import HomePage from './components/Landing.jsx/Landing';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Register from './components/GetStarted/Register';
import Footer from './components/Footer/Footer';
import Login from './components/Login/Login';
import CreateNote from './components/Note/CreateNote';
import Notes from './components/Note/Notes';
import ProtectedRoute from './components/Protect/ProtectedRoute';
import Fallback from './components/Fallback/Fallback';
import UpdateNote from './components/Note/UpdateNote';
function App() {


  return (
    <>
      <Router>
        <Header />
        <div className="min-h-screen flex flex-col items-center justify-center min-w-full">
          <main className="flex-grow sm:mt-[15vh] mt-[10vh] mb-10">
            <Routes>
              <Route path="/" element={<HomePage />} />
              <Route path="/register" element={<Register />} />
              <Route path="/login" element={<Login />} />
              <Route path="notes/update/:id" element={
                <ProtectedRoute>
                  <UpdateNote />
                </ProtectedRoute>
              } />
              <Route path="/create-note" element={
                <ProtectedRoute>
                  <CreateNote />
                </ProtectedRoute>
              } />
              <Route path="/notes" element={
                <ProtectedRoute>
                  <Notes />
                </ProtectedRoute>
              } />
              <Route path="*" element={<Fallback />} />
            </Routes>
          </main>
        </div>
        <Footer />
      </Router>
    </>
  )
}

export default App