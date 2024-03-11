import React, { useState } from 'react';

const UploadBox: React.FC = () => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file && file.type === 'application/pdf') {
      setSelectedFile(file);
    } else {
      alert('Please select a valid PDF file.');
    }
  };

  
  const handleSubmit = async () => {
    if (!selectedFile) return;
  
    const formData = new FormData();
    formData.append('pdf', selectedFile);
  
    try {
      const response = await fetch('http://localhost:3000/extract', {
        method: 'POST',
        body: formData,
      });
  
      if (response.ok) {
        const data = await response.json();
        console.log('Extracted data:', data);
        // Handle extracted reference and customer numbers here (if applicable)
      } else {
        console.error('Error processing PDF:', response.statusText);
      }
    } catch (error) {
      console.error('Error sending PDF:', error);
    }
  };
  
  
  return (
    <div>
      <input type="file" accept=".pdf" onChange={handleFileChange} />
      <button onClick={handleSubmit} disabled={!selectedFile}>
        Submit PDF
      </button>
    </div>
  );
};

export default UploadBox;
