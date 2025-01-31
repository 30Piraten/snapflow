// ... (Code to fetch the JSON from the Go backend. Example using fetch API)
const urlParams = new URLSearchParams(window.location.search);
const customerName = urlParams.get('customerName');
const orderID = urlParams.get('orderID');
const photoFileName = urlParams.get('photoFileName');

fetch(`/signed-url-info?customerName=${customerName}&orderID=${orderID}&photoFileName=${photoFileName}`)
// fetch('/signed-url-info?customerName=someName&orderID=someOrderId&photoFileName=photo.jpg') // Add query parameters
  .then(response => response.json())
  .then(signedURLInfo => {
      // 1. Retrieve Private Key from Secrets Manager
      AWS.config.region = 'YOUR_AWS_REGION'; // Set your AWS region

      const secretsManager = new AWS.SecretsManager();
      const secretName = 'YOUR_SECRET_NAME'; // The name of your secret

      secretsManager.getSecretValue({ SecretId: secretName }, (err, data) => {
          if (err) {
              console.error('Error retrieving secret:', err);
              return; // Stop execution if secret retrieval fails
          }

          const privateKey = data.SecretString;

          // 2. Generate Signed URL
          const signer = new AWS.CloudFront.Signer(signedURLInfo.keyPairId, privateKey);

          const url = signer.getSignedUrl({
              url: `https://${signedURLInfo.distributionDomain}/${signedURLInfo.objectKey}`,
              expires: signedURLInfo.expires,
              policy: signedURLInfo.policy
          });

          // 3. Use the URL (e.g., display image)
          const img = document.createElement('img');
          img.src = url;
          document.body.appendChild(img);
      });
  })
  .catch(error => {
      console.error('Error fetching signed URL info:', error);
  });