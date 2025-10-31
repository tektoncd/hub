#!//bin/bash
set -e -u -o pipefail

# Use writable directory for config.js (mounted via emptyDir volume)
CONFIG_DIR="${CONFIG_DIR:-/tmp/config}"
CONFIG_JS="${CONFIG_DIR}"/config.js

# Create directory if it doesn't exist
mkdir -p "${CONFIG_DIR}"

echo "Current ENV"
echo '----------------------------------------------'
if [ -f "$CONFIG_JS" ]; then
  cat $CONFIG_JS
  echo '..............................................'
  ls -l  $CONFIG_JS
else
  echo "Config file does not exist yet, will be created"
fi
echo "whoami: $(id)"
echo "Writing config to: $CONFIG_JS"
echo '----------------------------------------------'

cat <<-EOF > $CONFIG_JS
window.config = {
  API_URL: '$API_URL',
  API_VERSION: '$API_VERSION',
  AUTH_BASE_URL: '$AUTH_BASE_URL',
  REDIRECT_URI: '$REDIRECT_URI',
  CUSTOM_LOGO_MEDIA_TYPE: '$CUSTOM_LOGO_MEDIA_TYPE',
  CUSTOM_LOGO_BASE64_DATA: '$CUSTOM_LOGO_BASE64_DATA',
};


EOF

echo "Modified ENV"
cat $CONFIG_JS
echo '----------------------------------------------'


echo Starting Nginx
exec nginx -g 'daemon off;'

