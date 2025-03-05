import { useState } from 'react';
import LogoTitle from './LogoTitle';
import ToggleTheme from './../toggle-theme';
import NavBar from './NavBar';

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  return (
    <header className="fixed top-0 inset-x-0 z-50 bg-slate-500 shadow-lg">
      <div className="container px-6 py-4 mx-auto flex justify-between items-center">
        {/* Logo */}
        <div>
          <LogoTitle />
        </div>

        {/* Mobile Menu Button */}
        <div className="flex lg:hidden">
          <button
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            className="text-gray-200 hover:text-white focus:outline-none"
            aria-label="Toggle Menu"
          >
            {isMenuOpen ? (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="w-6 h-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            ) : (
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="w-6 h-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M4 8h16M4 16h16"
                />
              </svg>
            )}
          </button>
        </div>

        {/* Desktop Navigation */}
        <div className="hidden lg:flex items-center space-x-4">
          <NavBar setIsMenuOpen={setIsMenuOpen} />
          <ToggleTheme />
        </div>
      </div>

      {/* Mobile Navigation */}
      {isMenuOpen && (
        <div className="absolute inset-x-0 z-40 w-full px-6 py-4 bg-slate-500 lg:hidden">
          <NavBar setIsMenuOpen={setIsMenuOpen} />
          <div className="flex justify-center mt-4">
            <ToggleTheme />
          </div>
        </div>
      )}
    </header>
  );
};

export default Header;