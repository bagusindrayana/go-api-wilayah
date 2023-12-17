## SUMBER DATA
- https://github.com/guzfirdaus/Wilayah-Administrasi-Indonesia
- https://github.com/cahyadsn/wilayah
  
## INSTALASI/DEVELOPMENT
- import database `db_wilayah.sql`
- ubah konfigurasi database di `main.go` bagian `initDB()`

## API
DEMO : https://api-wilayah.up.railway.app

- `/info` : informasi jumlah provinsi, kabupaten/kota, kecamatan, dan desa
- `/provinsi` : list provinsi
- `/provinsi/{id}` : detail provinsi
- `/kota` : list kabupaten/kota
- `/provinsi/{id}/kabupaten` : list kabupaten/kota berdasarkan provinsi
- `/kota/{id}` : detail kabupaten/kota
- `/kecamatan` : list kecamatan
- `/kota/{id}/kecamatan` : list kecamatan berdasarkan kabupaten/kota
- `/kecamatan/{id}` : detail kecamatan
- `/desa` : list desa
- `/kecamatan/{id}/desa` : list desa berdasarkan kecamatan
- `/desa/{id}` : detail desa

### APIDOG
- https://4f6kku51q0.apidog.io