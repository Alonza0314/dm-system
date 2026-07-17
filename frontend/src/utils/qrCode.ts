import QRCode from 'qrcode'

export async function generateQrCodeDataUrl(url: string): Promise<string> {
  return QRCode.toDataURL(url, { width: 320, margin: 2 })
}

export function downloadDataUrl(dataUrl: string, filename: string) {
  const link = document.createElement('a')
  link.href = dataUrl
  link.download = filename
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
