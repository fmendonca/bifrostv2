import React from 'react';

function Spinner() {
  return (
    <div className="flex justify-center items-center h-full py-10">
      <div className="animate-spin rounded-full h-10 w-10 border-t-4 border-b-4 border-blue-500"></div>
    </div>
  );
}

export default Spinner;
