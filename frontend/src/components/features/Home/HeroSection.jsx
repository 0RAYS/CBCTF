import Button from '@/components/common/Button';
import { useNavigate } from 'react-router-dom';
import { useBranding } from '@/hooks/useBranding';

function HeroSection() {
  const navigate = useNavigate();
  const { home, homeLogo, homeLogoAlt } = useBranding();

  return (
    <div className="min-h-[560px] md:min-h-[700px] flex items-center px-4 md:px-8">
      <div className="page-container">
        <div className="grid grid-cols-1 md:grid-cols-[5fr_3fr] gap-12 md:gap-20 items-center">
          <div className="space-y-10">
            <div className="flex items-center gap-3">
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-geek-400 animate-pulse" />
                <span
                  className="w-1 h-1 rounded-full bg-geek-400/50 animate-pulse"
                  style={{ animationDelay: '0.3s' }}
                />
              </div>
              <span className="text-xs font-mono text-geek-400 tracking-[0.3em] uppercase">
                {home.hero.titlePrefix || 'CBCTF Platform'}
              </span>
              <div className="h-px flex-1 bg-geek-400/20 max-w-[120px]" />
            </div>

            <div>
              <h1 className="text-5xl sm:text-6xl md:text-7xl xl:text-8xl font-mono text-neutral-50 leading-[0.95] tracking-tight">
                {home.hero.titleHighlight}
              </h1>
              <p className="text-xl sm:text-2xl md:text-3xl font-mono text-neutral-500 mt-4 leading-snug tracking-tight">
                {home.hero.titleSuffix}
              </p>
            </div>

            <div className="flex items-center gap-4">
              <div className="w-8 h-[1px] bg-neutral-600" />
              <p className="text-neutral-400 text-sm leading-relaxed max-w-[50ch]">{home.hero.subtitle}</p>
            </div>

            <div className="flex flex-wrap gap-3">
              <Button variant="primary" size="lg" className="shadow-focus-strong" onClick={() => navigate('/games')}>
                {home.hero.primaryAction}
              </Button>
            </div>
          </div>

          <div className="relative hidden md:flex items-center justify-center">
            <div className="relative w-[280px] h-[280px] xl:w-[340px] xl:h-[340px]">
              <span className="absolute -top-3 -left-3 w-8 h-8 border-t-2 border-l-2 border-geek-400/70" />
              <span className="absolute -top-3 -right-3 w-8 h-8 border-t-2 border-r-2 border-geek-400/70" />
              <span className="absolute -bottom-3 -left-3 w-8 h-8 border-b-2 border-l-2 border-geek-400/70" />
              <span className="absolute -bottom-3 -right-3 w-8 h-8 border-b-2 border-r-2 border-geek-400/70" />
              <span className="absolute top-0 left-0 w-full h-full border border-neutral-700/30 rounded-sm" />

              <img
                alt={homeLogoAlt}
                src={homeLogo || './logo.png'}
                width="320"
                height="320"
                loading="lazy"
                decoding="async"
                className="w-full h-full object-contain p-8 opacity-85"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default HeroSection;
