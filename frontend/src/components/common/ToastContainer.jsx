import { motion, AnimatePresence } from 'motion/react';
import Toast from './Toast';

const POSITION_CLASSES = {
  'top-left': 'top-4 left-4',
  'top-center': 'top-4 left-1/2 -translate-x-1/2',
  'top-right': 'top-4 right-4',
  'bottom-left': 'bottom-4 left-4',
  'bottom-center': 'bottom-4 left-1/2 -translate-x-1/2',
  'bottom-right': 'bottom-4 right-4',
};

const ToastContainer = ({ position = 'top-right', toasts = [], removeToast }) => {
  const isTop = position.includes('top');
  const animations = {
    entry: { y: isTop ? -50 : 50, opacity: 0 },
    center: { y: 0, opacity: 1 },
    exit: { y: isTop ? -20 : 20, opacity: 0, transition: { duration: 0.2 } },
  };

  return (
    <div
      aria-live="polite"
      aria-atomic="false"
      aria-label="Notifications"
      className={`fixed z-[9999] flex flex-col ${POSITION_CLASSES[position] || POSITION_CLASSES['top-right']} pointer-events-none`}
      style={{
        gap: '0.75rem',
        maxWidth: 'calc(100vw - 2rem)',
        width: '420px',
        maxHeight: '100vh',
        overflowY: 'hidden',
      }}
    >
      <AnimatePresence mode="popLayout">
        {toasts.map((toast) => (
          <motion.div
            key={toast.id}
            layout
            initial={animations.entry}
            animate={animations.center}
            exit={animations.exit}
            transition={{
              type: 'spring',
              stiffness: 380,
              damping: 25,
              mass: 1,
            }}
            className="pointer-events-auto"
            style={{
              marginTop: '0.5rem',
              marginBottom: '0.5rem',
              filter: 'drop-shadow(0 5px 15px rgba(0, 0, 0, 0.2))',
              transformOrigin: isTop ? 'top' : 'bottom',
            }}
          >
            <Toast
              id={toast.id}
              title={toast.title}
              description={toast.description}
              color={toast.color}
              onClose={removeToast}
              hasCloseButton={toast.hasCloseButton !== false}
            />
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  );
};

export default ToastContainer;
