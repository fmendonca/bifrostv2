import React from 'react';

function Spinner() {
  return (
    <div className="flex justify-center items-center">
      <div className="animate-spin rounded-full h-12 w-12 border-t-4 border-b-4 border-bifrostBlue"></div>
    </div>
  );
}

export default Spinner;
