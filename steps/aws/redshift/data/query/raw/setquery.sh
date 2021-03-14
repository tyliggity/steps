export QUERY=$(cat <<EOF
SELECT *
FROM stv_sessions
WHERE user_name LIKE '%awsuser%';
EOF
)