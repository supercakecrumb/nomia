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
    <div className="min-h-screen flex flex-col bg-gray-50">
      {/* Header with Gradient Border */}
      <header className="bg-white shadow-md sticky top-0 z-50 border-b-2 border-transparent bg-gradient-to-r from-primary-500 via-secondary-500 to-accent-500" style={{ backgroundClip: 'border-box', backgroundOrigin: 'border-box', borderImageSlice: 1 }}>
        <div className="bg-white">
          <div className="container mx-auto px-4">
            <div className="flex items-center justify-between h-16">
              {/* Logo/Title with Icon */}
              <Link 
                to="/" 
                className="flex items-center gap-2 text-xl md:text-2xl font-bold text-gray-900 hover:text-primary-600 transition-colors group"
              >
                <div className="w-8 h-8 bg-gradient-to-br from-primary-600 to-accent-600 rounded-lg flex items-center justify-center shadow-md group-hover:shadow-lg transition-shadow">
                  <span className="text-white text-sm">✨</span>
                </div>
                <span className="bg-gradient-to-r from-primary-600 to-accent-600 bg-clip-text text-transparent">
                  {t('appName')}
                </span>
              </Link>

              {/* Navigation Links */}
              <nav className="hidden md:flex items-center gap-1">
                <Link
                  to="/"
                  className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                    isActive('/') 
                      ? 'bg-primary-100 text-primary-700' 
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
                  }`}
                >
                  {t('navigation.main')}
                </Link>
                <Link
                  to="/names"
                  className={`px-4 py-2 rounded-lg font-medium transition-all duration-200 ${
                    isActive('/names') 
                      ? 'bg-primary-100 text-primary-700' 
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'
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
            <nav className="md:hidden flex items-center gap-2 pb-3 border-t border-gray-200 mt-3 pt-3">
              <Link
                to="/"
                className={`flex-1 text-center py-2.5 px-4 rounded-lg font-medium transition-all duration-200 ${
                  isActive('/') 
                    ? 'bg-gradient-to-r from-primary-600 to-primary-700 text-white shadow-md' 
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {t('navigation.main')}
              </Link>
              <Link
                to="/names"
                className={`flex-1 text-center py-2.5 px-4 rounded-lg font-medium transition-all duration-200 ${
                  isActive('/names') 
                    ? 'bg-gradient-to-r from-primary-600 to-primary-700 text-white shadow-md' 
                    : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                }`}
              >
                {t('navigation.exploreNames')}
              </Link>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1">
        <Outlet />
      </main>

      {/* Enhanced Footer */}
      <footer className="bg-gradient-to-br from-gray-900 via-gray-800 to-gray-900 text-white py-12 mt-auto border-t-4 border-primary-600">
        <div className="container mx-auto px-4">
          <div className="grid md:grid-cols-3 gap-8 mb-8">
            {/* About Section */}
            <div>
              <div className="flex items-center gap-2 mb-4">
                <div className="w-8 h-8 bg-gradient-to-br from-primary-500 to-accent-500 rounded-lg flex items-center justify-center shadow-lg">
                  <span className="text-white text-sm">✨</span>
                </div>
                <h3 className="font-bold text-lg">{t('appName')}</h3>
              </div>
              <p className="text-gray-300 text-sm leading-relaxed">
                Supporting trans and nonbinary individuals in their journey to find affirming names through data-driven insights.
              </p>
            </div>

            {/* Quick Links */}
            <div>
              <h3 className="font-bold text-lg mb-4">Quick Links</h3>
              <ul className="space-y-2 text-sm">
                <li>
                  <Link 
                    to="/" 
                    className="text-gray-300 hover:text-white transition-colors hover:translate-x-1 inline-block"
                  >
                    → Home
                  </Link>
                </li>
                <li>
                  <Link 
                    to="/names" 
                    className="text-gray-300 hover:text-white transition-colors hover:translate-x-1 inline-block"
                  >
                    → Explore Names
                  </Link>
                </li>
              </ul>
            </div>

            {/* Mission */}
            <div>
              <h3 className="font-bold text-lg mb-4">Our Mission</h3>
              <p className="text-gray-300 text-sm leading-relaxed">
                Empowering identity exploration with comprehensive name statistics, gender-neutral insights, and inclusive design.
              </p>
            </div>
          </div>

          {/* Bottom Bar */}
          <div className="pt-8 border-t border-gray-700">
            <div className="flex flex-col md:flex-row justify-between items-center gap-4">
              <p className="text-sm text-gray-400">
                &copy; {new Date().getFullYear()} Affirm Name. All rights reserved.
              </p>
              <div className="flex items-center gap-4 text-sm text-gray-400">
                <span className="flex items-center gap-2">
                  Made with 
                  <span className="text-accent-400 animate-pulse">❤️</span>
                  for the community
                </span>
              </div>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}