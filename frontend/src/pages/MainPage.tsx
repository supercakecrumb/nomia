import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

export default function MainPage() {
  const { t } = useTranslation('pages');

  return (
    <div className="min-h-screen">
      {/* Hero Section with Gradient Background */}
      <section className="relative overflow-hidden bg-gradient-to-br from-primary-50 via-secondary-50 to-accent-50">
        <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNhODU1ZjciIGZpbGwtb3BhY2l0eT0iMC4wNSI+PHBhdGggZD0iTTM2IDEzNGg3NnYxNkg0MVYxMzRoLTY3djE2aDE2em0wLTE2aDE2di0xNmg0NXYxNkg0MVYxMThoLTZ2MTZoLTZ2LTE2SDN2MTZoNnYxNmg2di0xNmg2ek00MVwyNTBoNzZ2LTE2SDQxdjE2em02LTc2aDc2djE2SDUxdi0xNmgtNnYxNmgtNnYtMTZIM3YxNmg2djE2aDZ2LTE2aDZ6TTQxIDIzNGg3NnYxNkg0MXYtMTZ6Ii8+PC9nPjwvZz48L3N2Zz4=')] opacity-40"></div>
        
        <div className="container mx-auto px-4 py-20 md:py-32 relative">
          <div className="max-w-4xl mx-auto text-center">
            <div className="animate-fade-in">
              <h1 className="text-5xl md:text-7xl font-bold text-gray-900 mb-6 tracking-tight">
                <span className="block">{t('main.hero.title')}</span>
              </h1>
              <p className="text-xl md:text-2xl text-gray-700 mb-12 leading-relaxed max-w-3xl mx-auto">
                {t('main.hero.tagline')}
              </p>
              
              {/* CTA Button with Gradient and Animation */}
              <Link
                to="/names"
                className="inline-flex items-center gap-2 bg-gradient-to-r from-primary-600 to-primary-700 hover:from-primary-700 hover:to-primary-800 text-white font-semibold px-8 py-4 rounded-xl text-lg transition-all duration-300 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95 group"
              >
                {t('main.cta.button')}
                <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                </svg>
              </Link>
            </div>
          </div>
        </div>

        {/* Decorative Wave Divider */}
        <div className="absolute bottom-0 left-0 right-0">
          <svg className="w-full h-16 text-white" preserveAspectRatio="none" viewBox="0 0 1200 120" xmlns="http://www.w3.org/2000/svg">
            <path d="M321.39,56.44c58-10.79,114.16-30.13,172-41.86,82.39-16.72,168.19-17.73,250.45-.39C823.78,31,906.67,72,985.66,92.83c70.05,18.48,146.53,26.09,214.34,3V0H0V27.35A600.21,600.21,0,0,0,321.39,56.44Z" fill="currentColor"></path>
          </svg>
        </div>
      </section>

      {/* Mission Section */}
      <section className="container mx-auto px-4 py-16 md:py-24">
        <div className="max-w-4xl mx-auto">
          <div className="bg-white rounded-2xl shadow-medium hover:shadow-strong transition-shadow duration-300 p-8 md:p-12 border border-gray-100">
            <div className="flex items-start gap-4 mb-6">
              <div className="flex-shrink-0">
                <div className="w-14 h-14 bg-gradient-to-br from-primary-500 to-primary-600 rounded-xl flex items-center justify-center shadow-lg">
                  <svg className="w-7 h-7 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                  </svg>
                </div>
              </div>
              <div className="flex-1">
                <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
                  {t('main.mission.title')}
                </h2>
                <p className="text-lg text-gray-700 leading-relaxed">
                  {t('main.mission.description')}
                </p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="bg-gray-50 py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="max-w-6xl mx-auto">
            <div className="text-center mb-16">
              <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
                {t('main.features.title')}
              </h2>
              <p className="text-lg text-gray-600 max-w-2xl mx-auto">
                Powerful tools to help you find the perfect name
              </p>
            </div>

            <div className="grid md:grid-cols-2 gap-6">
              {/* Feature 1 */}
              <div className="group bg-white rounded-2xl shadow-soft hover:shadow-strong transition-all duration-300 p-8 border border-gray-100 hover:border-primary-200">
                <div className="flex items-start gap-4">
                  <div className="flex-shrink-0">
                    <div className="w-14 h-14 bg-gradient-to-br from-primary-100 to-primary-200 rounded-xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform duration-300">
                      <span className="text-3xl">🔍</span>
                    </div>
                  </div>
                  <div className="flex-1">
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">
                      Explore Names
                    </h3>
                    <p className="text-gray-700">
                      {t('main.features.explore')}
                    </p>
                  </div>
                </div>
              </div>

              {/* Feature 2 */}
              <div className="group bg-white rounded-2xl shadow-soft hover:shadow-strong transition-all duration-300 p-8 border border-gray-100 hover:border-secondary-200">
                <div className="flex items-start gap-4">
                  <div className="flex-shrink-0">
                    <div className="w-14 h-14 bg-gradient-to-br from-secondary-100 to-secondary-200 rounded-xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform duration-300">
                      <span className="text-3xl">🎯</span>
                    </div>
                  </div>
                  <div className="flex-1">
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">
                      Advanced Filters
                    </h3>
                    <p className="text-gray-700">
                      {t('main.features.filter')}
                    </p>
                  </div>
                </div>
              </div>

              {/* Feature 3 */}
              <div className="group bg-white rounded-2xl shadow-soft hover:shadow-strong transition-all duration-300 p-8 border border-gray-100 hover:border-accent-200">
                <div className="flex items-start gap-4">
                  <div className="flex-shrink-0">
                    <div className="w-14 h-14 bg-gradient-to-br from-accent-100 to-accent-200 rounded-xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform duration-300">
                      <span className="text-3xl">📊</span>
                    </div>
                  </div>
                  <div className="flex-1">
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">
                      Data Trends
                    </h3>
                    <p className="text-gray-700">
                      {t('main.features.trends')}
                    </p>
                  </div>
                </div>
              </div>

              {/* Feature 4 */}
              <div className="group bg-white rounded-2xl shadow-soft hover:shadow-strong transition-all duration-300 p-8 border border-gray-100 hover:border-primary-200">
                <div className="flex items-start gap-4">
                  <div className="flex-shrink-0">
                    <div className="w-14 h-14 bg-gradient-to-br from-primary-100 to-accent-200 rounded-xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform duration-300">
                      <span className="text-3xl">✨</span>
                    </div>
                  </div>
                  <div className="flex-1">
                    <h3 className="text-xl font-semibold text-gray-900 mb-2">
                      Affirming Choice
                    </h3>
                    <p className="text-gray-700">
                      {t('main.features.affirming')}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Call to Action Section */}
      <section className="container mx-auto px-4 py-16 md:py-24">
        <div className="max-w-4xl mx-auto">
          <div className="relative overflow-hidden bg-gradient-to-r from-primary-600 to-accent-600 rounded-3xl shadow-strong p-12 md:p-16 text-center">
            {/* Decorative Elements */}
            <div className="absolute top-0 right-0 -mt-4 -mr-4 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
            <div className="absolute bottom-0 left-0 -mb-4 -ml-4 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
            
            <div className="relative">
              <h2 className="text-3xl md:text-4xl font-bold text-white mb-4">
                Ready to Start Your Journey?
              </h2>
              <p className="text-xl text-white/90 mb-8 max-w-2xl mx-auto">
                {t('main.cta.subtitle')}
              </p>
              <Link
                to="/names"
                className="inline-flex items-center gap-2 bg-white text-primary-700 font-semibold px-8 py-4 rounded-xl text-lg transition-all duration-300 shadow-lg hover:shadow-2xl hover:scale-105 active:scale-95 group"
              >
                {t('main.cta.button')}
                <svg className="w-5 h-5 transform group-hover:translate-x-1 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                </svg>
              </Link>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}