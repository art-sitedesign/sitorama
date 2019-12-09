echo "<h4>Postgres checker:</h4>";

$connection = new PDO('pgsql:host={{.DBHost}};dbname={{.DBName}}', '{{.DBUser}}', '{{.DBPass}}');

$result = $connection->query("select 'OK' AS result");

foreach ($result as $row) {
    echo "DB query result: ", $row['result'];
}
