import React, { useState } from 'react';
export const VersionContext = React.createContext({
  version: '',
  setVersion: undefined
});

export default function VersionContextProvider({ children }) {
  const [version, setVersion] = useState('');
  const obj = {
    version,
    setVersion
  };

  return (
    <VersionContext.Provider value={obj}>{children}</VersionContext.Provider>
  );
}
