import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

export default function MainPage() {
  const { t } = useTranslation('pages');

  return (
    <div className="min-h-screen bg-gradient-to-b from-blue-50 to-white">
      {/* Hero Section */}
      <section className="container mx-auto px-4 py-16 md:py-24">
        <div className="max-w-4xl mx-auto text-center">
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 mb-6">
            {t('main.hero.title')}
          </h1>
          <p className="text-xl md:text-2xl text-gray-600 mb-12">
            {t('main.hero.tagline')}
          </p>
        </div>
      </section>

      {/* Mission Section */}
      <section className="container mx-auto px-4 py-12 md:py-16">
        <div className="max-w-4xl mx-auto bg-white rounded-lg shadow-lg p-8 md:p-12">
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            {t('main.mission.title')}
          </h2>
          <p className="text-lg text-gray-700 leading-relaxed">
            {t('main.mission.description')}
          </p>
        </div>
      </section>

      {/* Features Section */}
      <section className="container mx-auto px-4 py-12 md:py-16">
        <div className="max-w-4xl mx-auto">
          <h2 className="text-3xl md:text-4xl font-bold text-gray-900 mb-12 text-center">
            {t('main.features.title')}
          </h2>
          <div className="grid md:grid-cols-2 gap-6">
            <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="flex items-start gap-4">
                <div className="flex-shrink-0">
                  <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">🔍</span>
                  </div>
                </div>
                <div>
                  <p className="text-gray-700">
                    {t('main.features.explore')}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="flex items-start gap-4">
                <div className="flex-shrink-0">
                  <div className="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">🎯</span>
                  </div>
                </div>
                <div>
                  <p className="text-gray-700">
                    {t('main.features.filter')}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="flex items-start gap-4">
                <div className="flex-shrink-0">
                  <div className="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">📊</span>
                  </div>
                </div>
                <div>
                  <p className="text-gray-700">
                    {t('main.features.trends')}
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="flex items-start gap-4">
                <div className="flex-shrink-0">
                  <div className="w-12 h-12 bg-pink-100 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">✨</span>
                  </div>
                </div>
                <div>
                  <p className="text-gray-700">
                    {t('main.features.affirming')}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Call to Action */}
      <section className="container mx-auto px-4 py-12 md:py-16 mb-12">
        <div className="max-w-4xl mx-auto text-center">
          <p className="text-xl text-gray-600 mb-8">
            {t('main.cta.subtitle')}
          </p>
          <Link
            to="/names"
            className="inline-block bg-blue-600 hover:bg-blue-700 text-white font-semibold px-8 py-4 rounded-lg text-lg transition-colors shadow-lg hover:shadow-xl"
          >
            {t('main.cta.button')}
          </Link>
        </div>
      </section>
    </div>
  );
}