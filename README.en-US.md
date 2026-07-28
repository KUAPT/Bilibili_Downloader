# Bilibili_Downloader  
> This project is a simple Bilibili video downloader, featuring a built-in transcoding function based on embedded ffmpeg to output directly playable mp4 format videos.   
> This project is for development and learning purposes only and is not intended to infringe upon the rights of any party. If you believe this project infringes upon your legal rights, please contact the developer, and it will be handled as soon as possible.  
> This project is open-sourced under the GPLv3 license.   
> License: [GNU General Public License v3.0](https://github.com/KUAPT/Bilibili_Downloader/blob/main/README.md)  

## How to Use:
1. Go to [Release](https://github.com/KUAPT/Bilibili_Downloader/releases) to download the latest version (only the latest version is maintained and bug-fixed).  
2. Place the downloaded executable file in the directory of your choice (it is recommended to place it in an **empty directory** that is easy to find).  
3. Double-click to run. ~~Basically, just run it.~~    

> ps: If the program fails to run normally on the first attempt (including but not limited to crashing, which may be due to insufficient program permissions, try running with administrator privileges or manually granting the necessary permissions).  

4. The first run will generate necessary files and directories where the program is located. **Please do not delete them unless necessary** to avoid affecting the user experience.  
5. When not logged in (which is required to obtain Cookies for downloading videos in 1080p and above), a login QR code will be provided. Please use the Bilibili client (app) to scan and log in (~~the scan button should be on the personal homepage~~). Theoretically, this operation only needs to be done once, as the Cookie information will be persisted to simplify the subsequent process.  

> ps: If you find that video downloads fail or you cannot download videos higher than 720p, try deleting the `config` directory and restarting the program to log in again to rule out potential Cookie expiration issues.  

6. Enter the video BV number (~~found in the video details in the app or in the address bar on the web~~) and wait for the download.  
7. Once the download is complete, the `mp4` video will be output to the `Download` directory within the program's location.  

## Tips:
1. If the login QR code displays abnormally when running with administrator privileges, try running it directly as a standard user.  
2. This project **will not upload any personal information**. Information such as Cookies will be stored locally in the `config` directory. The developer is not responsible for privacy leaks caused by improper personal management.  
3. The automatic update of this project depends on Github. If Github is inaccessible due to network ~~environment~~ issues, the automatic update will fail.  

## ~~Pie in the Sky~~ Future Development Plan (Todo):
- [x] Quality selection function  
- [x] Multi-part (P) selection/continuous download function  
- [x] Auto-detect and update  
- [ ] Compatibility for both BV numbers and URLs  
- [ ] Asynchronous processing for continuous downloads  
- [ ] Custom output directory function  
- [ ] Custom output format function  
- [ ] GUI graphical interface (the functions are so ~~simple~~ pure, there might not be one in the short term :) )
