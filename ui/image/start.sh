#!//bin/bash
set -e -u -o pipefail

CONFIG_JS=/usr/share/nginx/html/config.js
echo "Current ENV"
echo '----------------------------------------------'
cat $CONFIG_JS
echo '..............................................'
ls -l  $CONFIG_JS
echo "whoami: $(id)"
echo '----------------------------------------------'

cat <<-EOF > $CONFIG_JS
window.config = {
  API_URL: '$API_URL',
  API_VERSION: '$API_VERSION',
  AUTH_BASE_URL: '$AUTH_BASE_URL',
  REDIRECT_URI: '$REDIRECT_URI',
};


EOF

echo "Modified ENV"
cat $CONFIG_JS
echo '----------------------------------------------'


echo Starting Nginx
exec nginx -g 'daemon off;'

