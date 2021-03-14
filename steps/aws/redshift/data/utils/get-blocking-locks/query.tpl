SELECT
    a.pid,
    a.xid,
    a.pidlist,
    a.username,
    a.block_sec,
    a.max_sec_blocking,
    a.num_blocking,
    b.querytxt
FROM admin.v_get_blocking_locks a
LEFT JOIN stl_query b on b.xid=a.xid
WHERE num_blocking>0 {{if .Conditions}}{{.Conditions}}{{end}}
GROUP BY 1,2,3,4,5,6,7,8
ORDER BY a.block_sec asc