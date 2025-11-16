import { useTranslation } from 'react-i18next';
import { Link, Outlet, useLocation } from 'react-router-dom';
import LanguageSwitcher from '../LanguageSwitcher';

export default function AppLayout() {
  const { t } = useTranslation('common');
  const location = useLocation();

  const isActive = (path: string) => {
    return location.pathname === path;
  };

  return (
    <div className="min-h-screen flex flex-col">
      {/* Header */}
      <header className="bg-white shadow-md sticky top-0 z-50">
        <div className="container mx-auto px-4">
          <div className="flex items-center justify-between h-16">
            {/* Logo/Title */}
            <Link 
              to="/" 
              className="text-xl md:text-2xl font-bold text-gray-900 hover:text-blue-600 transition-colors"
            >
              {t('appName')}
            </Link>

            {/* Navigation Links */}
            <nav className="hidden md:flex items-center gap-6">
              <Link
                to="/"
                className={`font-medium transition-colors ${
                  isActive('/') 
                    ? 'text-blue-600' 
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                {t('navigation.main')}
              </Link>
              <Link
                to="/names"
                className={`font-medium transition-colors ${
                  isActive('/names') 
                    ? 'text-blue-600' 
                    : 'text-gray-600 hover:text-gray-900'
                }`}
              >
                {t('navigation.exploreNames')}
              </Link>
            </nav>

            {/* Language Switcher */}
            <div className="flex items-center gap-4">
              <LanguageSwitcher />
            </div>
          </div>

          {/* Mobile Navigation */}
          <nav className="md:hidden flex items-center gap-4 pb-3 border-t border-gray-200 mt-3 pt-3">
            <Link
              to="/"
              className={`flex-1 text-center py-2 px-4 rounded-lg font-medium transition-colors ${
                isActive('/') 
                  ? 'bg-blue-600 text-white' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {t('navigation.main')}
            </Link>
            <Link
              to="/names"
              className={`flex-1 text-center py-2 px-4 rounded-lg font-medium transition-colors ${
                isActive('/names') 
                  ? 'bg-blue-600 text-white' 
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
            >
              {t('navigation.exploreNames')}
            </Link>
          </nav>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="bg-gray-800 text-white py-8 mt-auto">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-3 gap-8">
            <div>
              <h3 className="font-bold text-lg mb-3">{t('appName')}</h3>
              <p className="text-gray-300 text-sm">
                Supporting trans and nonbinary individuals in their journey to find affirming names.
              </p>
            </div>
            <div>
              <h3 className="font-bold text-lg mb-3">Quick Links</h3>
              <ul className="space-y-2 text-sm">
                <li>
                  <Link to="/" className="text-gray-300 hover:text-white transition-colors">
                    Home
                  </Link>
                </li>
                <li>
                  <Link to="/names" className="text-gray-300 hover:text-white transition-colors">
                    Explore Names
                  </Link>
                </li>
              </ul>
            </div>
            <div>
              <h3 className="font-bold text-lg mb-3">About</h3>
              <p className="text-gray-300 text-sm">
                This project provides data-driven insights into name trends to help you make informed decisions about your identity.
              </p>
            </div>
          </div>
          <div className="mt-8 pt-8 border-t border-gray-700 text-center text-sm text-gray-400">
            <p>&copy; {new Date().getFullYear()} Affirm Name. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}