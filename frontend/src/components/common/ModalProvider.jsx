import { createContext, useContext, useState } from 'react';

const ModalPortalContext = createContext(null);

// eslint-disable-next-line react-refresh/only-export-components
export function useModalPortal() {
  return useContext(ModalPortalContext);
}

export default function ModalProvider({ children }) {
  const [container, setContainer] = useState(null);

  return (
    <ModalPortalContext.Provider value={container}>
      {children}
      <div ref={setContainer} />
    </ModalPortalContext.Provider>
  );
}
