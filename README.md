Source-code for the fundraising-website of our photovoltaic-system

# TODO
- move prices from list to popup?
- move list of valid mids and regex out of source-code
- create mail-template-file
- add second certificate for empty names
- improve certificate-email text and header
- certificate: use better formatted names
- simplify uploading non-git files to server (fonts, elements, inkscape, certificate-template) rsync-script?
- johannes-logo?
- deploy: restart backend
- change workflow to only run on tag
- check npm package versions / audit
- popup-close-x isn't in right corner

rsync --delete -a client/src/public/external z2glr@z2glr.uber.space:~/pv-spenden/client/src/public/external