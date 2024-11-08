Source-code for the fundraising-website of our photovoltaic-system

# TODO
- simplify uploading non-git files to server (fonts, elements, inkscape, certificate-template) rsync-script?
- change workflow to only run on tag
- check npm package versions / audit
- refresh jwt on login

rsync --delete -a client/src/public/external z2glr@z2glr.uber.space:~/pv-spenden/client/src/public/external
--files-from=FILE

# files to manually copy
- images
- elements
- valid-elements-list
- e-mail-templates
- certificate templates