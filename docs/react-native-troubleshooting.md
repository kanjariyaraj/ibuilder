# React Native Troubleshooting

## Common Issues

### Node.js not found

```
builder rn doctor
```

Check Node.js is installed and in PATH:
```
node --version
```

### Metro won't start (port conflict)

```
lsof -i :8081
kill -9 <PID>
builder rn metro start
```

### Device not detected

Ensure MobAI is connected:
```
builder mobai status
builder device list
```

### Dependencies not installed

```
npm install
```

### Build fails on Windows

If using Windows, prefer WSL2 for iOS development:
```
builder doctor
```

### MobAI disconnection

Builder auto-reconnects on disconnect.
Check connection:
```
builder mobai status
```

### Log streaming not showing output

Ensure the app is running on the device:
```
builder rn logs --stream
```
