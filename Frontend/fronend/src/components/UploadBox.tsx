// import React, { useState } from 'react';
// import axios from 'axios';
// import DisplayArea from './DisplayArea';

// interface ExtractedData {
//   referenceNumber: string;
//   customerNumber: string;
// }

// const UploadBox: React.FC = () => {
//   const [selectedFile, setSelectedFile] = useState<File | null>(null);
//   const [isLoading, setIsLoading] = useState(false);
//   const [error, setError] = useState<string | null>(null);
//   const [data, setData] = useState<ExtractedData | null>(null);

//   const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
//     const file = event.target.files?.[0];
//     setSelectedFile(file);
//   };

//   const handleSubmit = async () => {
//     if (!selectedFile) return;

//     setIsLoading(true);
//     setError(null);

//     try {
//       const formData = new FormData();
//       formData.append('pdf', selectedFile);

//       // Replace with your backend API endpoint when ready
//       const response = await axios.post('/api/extract-data', formData, {
//         headers: { 'Content-Type': 'multipart/form-data' },
//       });

//       setData(response.data);
//     } catch (error) {
//       setError(error.message);
//     } finally {
//       setIsLoading(false);
//     }
//   };

//   return (
//     <div>
//       <input type="file" onChange={handleFileChange} />
//       <button onClick={handleSubmit} disabled={isLoading}>
//         {isLoading ? 'Extracting...' : 'Extract Data'}
//       </button>
//       {error && <p style={{ color: 'red' }}>{error}</p>}
//       {data && <DisplayArea data={data} />}
//     </div>
//   );
// };

// export default UploadBox;
// -------------------2nd----------------
// import React, { useState } from 'react';

// const UploadBox: React.FC = () => {
//   const [selectedFile, setSelectedFile] = useState<File | null>(null);

//   const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
//     const file = event.target.files?.[0];
//     if (file && file.type === 'application/pdf') { // Ensure only PDF files are selected
//       setSelectedFile(file);
//     } else {
//       alert('Please select a valid PDF file.');
//     }
//   };

//   return (
//     <div>
//       <input type="file" accept=".pdf" onChange={handleFileChange} />
//       <p>Selected File: {selectedFile?.name || 'No file selected'}</p>
//       {/* Comment out or remove the submit button and any backend-related actions */}
//       {/* <button onClick={handleSubmit} disabled={!selectedFile}>Extract Data</button> */}
//     </div>
//   );
// };

// export default UploadBox;
// ------------------------------

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

  // const handleSubmit = () => {
  //   if (!selectedFile) return;

  //   // **DISCLAIMER:** This doesn't actually store the file permanently.
  //   // It's only for demonstration purposes to simulate form submission.
  //   const reader = new FileReader();
  //   reader.readAsDataURL(selectedFile);
  //   reader.onload = () => {
  //     const fileData = reader.result; // This is a temporary variable
  //     console.log('File data:', fileData); // Simulate processing
  //   };
  // };
  // const handleSubmit = async () => {
  //   if (!selectedFile) return;
  
  //   const formData = new FormData();
  //   formData.append('pdf', selectedFile);
  
  //   try {
  //     const response = await fetch('/extract', {
  //       method: 'POST',
  //       body: formData,
  //     });
  
  //     if (response.ok) {
  //       const data = await response.json();
  //       console.log('Extracted data:', data);
  //       // Handle extracted reference and customer numbers here (if applicable)
  //     } else {
  //       console.error('Error processing PDF:', response.statusText);
  //     }
  //   } catch (error) {
  //     console.error('Error sending PDF:', error);
  //   }
  // };
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
