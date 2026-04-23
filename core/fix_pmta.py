import os

path = r"C:\Users\ayoub\Documents\BillionMail-cold\core\internal\cmd\cmd.go"

with open(path, 'rb') as f:
    data = f.read()

old = b'''						// check if the request is in the excluded URIs
						if _, ok := excludesURIs[r.URL.Path]; ok {
							return
						}

						if r.IsFileRequest() {'''

new = b'''						// check if the request is in the excluded URIs
						if _, ok := excludesURIs[r.URL.Path]; ok {
							return
						}

						// Allow email tracking paths (pixel + click redirects)
						if strings.HasPrefix(r.URL.Path, "/pmta/") {
							return
						}

						// Allow public landing pages
						if strings.HasPrefix(r.URL.Path, "/landing/") {
							return
						}

						if r.IsFileRequest() {'''

if old in data:
    data = data.replace(old, new, 1)
    with open(path, 'wb') as f:
        f.write(data)
    print("SUCCESS: SafePath bypass for /pmta/ and /landing/ added")
elif new in data:
    print("ALREADY PATCHED")
else:
    print("ERROR: pattern not found")
    # Try to find nearby content for debugging
    idx = data.find(b'excludesURIs[r.URL.Path]')
    if idx >= 0:
        print(f"Found excludesURIs at byte {idx}")
        print(repr(data[idx-50:idx+200]))
