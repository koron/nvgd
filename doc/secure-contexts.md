# Secure Contexts

Certain NVGD features (e.g., OPFS) require a [secure context][mdnsc].
A secure context refers to a connection established via localhost, a TLS (HTTPS) connection, or a connection to an origin explicitly marked as secure through browser flags.
NVGD is designed with the assumption that it will be operated within a private Local Area Network (LAN).
In such environments, obtaining a formal TLS certificate from a Certificate Authority (CA) can be challenging.
Therefore, using TLS with a self-signed certificate is a common approach.

This document explains two methods:
1. How to configure your browser to treat a specific origin as a secure context, even if it's not served over HTTPS.
2. How to generate a self-signed certificate and configure NVGD to use it for providing TLS.

[mdnsc]:https://developer.mozilla.org/en-US/docs/Web/Security/Secure_Contexts

## Marking an Insecure Origin as Secure via Browser Flags

This section describes how to instruct your browser to treat an insecure origin as secure.

1.  In your browser, navigate to `chrome://flags/#unsafely-treat-insecure-origin-as-secure`.

    *   This flag is specific to Chrome-based browsers.
    *   For Microsoft Edge, use `edge://flags/#unsafely-treat-insecure-origin-as-secure`.
    *   As of June 22, 2025, Firefox does not offer an equivalent setting.

2.  In the input field, enter the URL(s) you want to treat as secure.

    *   Ensure you include the scheme (e.g., `http://`) and the port number.
    *   To specify multiple URLs, separate them with commas (`,`).

    Example:

    ```
    http://192.168.0.100:9280,http://192.168.0.101:9280,http://dev.mydomain.org:9280
    ```

3. Enable the flag and relaunch your browser when prompted.

## Creating and Using a Self-Signed Certificate with NVGD

This section outlines the process for creating a self-signed certificate and configuring NVGD to use it for TLS. Using a self-signed certificate is suitable for development, testing, or internal LAN environments where obtaining a CA-signed certificate is not feasible.

**Note:** Clients (browsers, command-line tools) will need to be configured to trust this self-signed certificate, as it won't be signed by a recognized Certificate Authority (CA).

### 1. Generate a Self-Signed Certificate and Private Key

You can use OpenSSL to generate a private key and a self-signed certificate.

**a. Generate a private key:**

```bash
openssl genpkey -algorithm RSA -out server.key
```
This command creates a new RSA private key and saves it to `server.key`.

**b. Generate a Certificate Signing Request (CSR):**
While not strictly necessary for a self-signed certificate that you sign yourself, creating a CSR is good practice if you want to include specific information in your certificate.

```bash
openssl req -new -key server.key -out server.csr
```
You will be prompted to enter information such as country name, organization name, and common name. For the "Common Name (e.g. server FQDN or YOUR name)", use the hostname or IP address that clients will use to access NVGD (e.g., `localhost`, `192.168.1.10`, `nvgd.local`).

**c. Generate the self-signed certificate:**

```bash
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt
```
This command takes the CSR (`server.csr`) and the private key (`server.key`), signs the request with the key, and outputs a self-signed certificate (`server.crt`) valid for 365 days.

You should now have two essential files:
*   `server.key`: Your private key. **Keep this file secure.**
*   `server.crt`: Your self-signed certificate.

### 2. Configure and Start NVGD in TLS Mode

NVGD needs to be configured to use the generated `server.crt` (certificate file) and `server.key` (key file).

The exact method for configuring TLS in NVGD will depend on its command-line flags or configuration file format. Typically, you would specify paths to the certificate and key files.

Example (hypothetical command-line flags):

```bash
nvgd --tls-cert /path/to/your/server.crt --tls-key /path/to/your/server.key
```

Or, if NVGD uses a configuration file (e.g., `config.yml`):

```yaml
# Example config.yml snippet
server:
  tls:
    enabled: true
    cert_file: "/path/to/your/server.crt"
    key_file: "/path/to/your/server.key"
  # ... other server configurations
```

**Consult the NVGD documentation for the precise command-line arguments or configuration file options related to TLS.**

Once configured, start NVGD. It should now be serving HTTPS requests on its configured port.

### 3. Configure Clients to Trust the Self-Signed Certificate

When you try to access NVGD over HTTPS using a browser or other client, you will likely see a security warning because the certificate is self-signed and not trusted by a known CA.

**a. For Browsers:**

Most browsers will allow you to add an exception for the self-signed certificate. The steps vary by browser:
*   You might see an "Advanced" button on the warning page, leading to an option like "Proceed to [hostname] (unsafe)" or "Accept the Risk and Continue."
*   Alternatively, you might need to import `server.crt` into your browser's certificate trust store. Search for your browser's documentation on "importing trusted root certificates" or "managing SSL certificates."

**b. For Command-Line Tools (e.g., `curl`):**

*   **Insecure Option (for testing only):** Many tools have an "insecure" flag (e.g., `curl -k` or `curl --insecure`) that bypasses certificate validation. **This is not recommended for production or sensitive environments.**
*   **Specify CA Certificate:** A better approach is to tell the client to trust your specific self-signed certificate. For `curl`, you can use the `--cacert` option:
    ```bash
    curl --cacert /path/to/your/server.crt https://your-nvgd-host:port/
    ```
    Some tools or libraries might require the certificate to be added to the system's trust store.

**Important Considerations:**

*   **Security:** Self-signed certificates provide encryption but do not offer the same level of trust and identity verification as CA-signed certificates. They are vulnerable to man-in-the-middle attacks if the certificate fingerprint is not verified out-of-band.
*   **Distribution:** If multiple clients need to access NVGD, each client will need to be configured to trust the certificate.
*   **Expiration:** Self-signed certificates expire. Remember the validity period you set (e.g., 365 days) and plan to renew the certificate before it expires.
*   **Alternative for Internal Networks: Private CA:** For more robust internal deployments, consider setting up your own private Certificate Authority (CA). You would then issue certificates signed by your private CA and install the private CA's root certificate on all client machines. This avoids having to trust individual self-signed certificates on each client.
