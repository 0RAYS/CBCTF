import { useState } from 'react';

const SIZE_MAP = {
  xs: 32,
  sm: 40,
  md: 48,
  lg: 64,
  xl: 96,
};

function Avatar({ src, name = '', size = 'md', shape = 'rounded', className = '' }) {
  const [imgError, setImgError] = useState(false);
  const [prevSrc, setPrevSrc] = useState(src);

  if (prevSrc !== src) {
    setPrevSrc(src);
    setImgError(false);
  }

  const px = typeof size === 'number' ? size : SIZE_MAP[size] || 48;
  const radius = shape === 'circle' ? '9999px' : '0.5rem';
  const showFallback = !src || imgError;
  const initial = name ? name[0].toUpperCase() : '?';
  const fontSize = Math.max(12, Math.round(px * 0.4));

  return (
    <div className={`shrink-0 overflow-hidden ${className}`} style={{ width: px, height: px, borderRadius: radius }}>
      {showFallback ? (
        <div className="w-full h-full bg-neutral-700 flex items-center justify-center" style={{ fontSize }}>
          <span className="font-mono text-neutral-300 leading-none">{initial}</span>
        </div>
      ) : (
        <img src={src} alt={name} loading="lazy" className="w-full h-full object-cover" onError={() => setImgError(true)} />
      )}
    </div>
  );
}

export default Avatar;
