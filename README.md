```	
# Play mp4 file
ffplay mini_bbb.mp4

# Check hash of mp4 file
sha256sum mini_bbb.mp4

# Mux to mpegts
ffmpeg -i mini_bbb.mp4 -c copy -bsf h264_mp4toannexb mini_bbb.ts

# Play mpegts file
ffplay mini_bbb.ts

# Check hash of mpegts file
sha256sum mini_bbb.ts

# Show frame hashes of mp4 file
ffmpeg -i data/mini_bbb.mp4 -f framehash -hash sha256 -

# Show frame hashes of mpegts file
ffmpeg -i data/mini_bbb.ts -f framehash -hash sha256 -

# Generate key pair and store private key in mykey
go run main.go gen-key mykey

# Generate root RMID for mp4 file
go run main.go root-rmid data/mini_bbb.mp4

# Generate root RMID for mpegts file
go run main.go root-rmid data/mini_bbb.ts

# Sign frame hashes of mp4 file using private key in mykey and write to sigs.txt
go run main.go sign-video data/mini_bbb.mp4 mykey > sigs.txt

# Verify signatures of frame hashes of mp4 file against public key
go run main.go sign-video data/mini_bbb.ts <PUBLIC_KEY> sigs.txt

# Verify signatures of frame hashes of mpegts file against public key
go run main.go sign-video data/mini_bbb.ts <PUBLIC_KEY> sigs.txt
```