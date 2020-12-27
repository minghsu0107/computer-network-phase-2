# Video Streamer 
## Usage
### Requirements
- go
- npm
- bash
### Setup
Build React frontend:
```bash
cd ui
npm install
npm run build
```
Run backend:
```bash
./scripts/run.sh serve
```
### Create a configuration file
Example configuration format is as followed:
```
{
	"folders": [
		{
			"path": "/data/Videos",
			"title": "My Videos"
		},
		{
			"path": "/data/Others",
			"title": "My Others"
		}
	]
}
```