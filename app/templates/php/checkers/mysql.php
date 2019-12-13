echo "<h4>MySQL checker:</h4>";

$connection = new PDO('mysql:host={{.DBHost}};dbname={{.DBName}}', '{{.DBUser}}', '{{.DBPass}}');

$result = $connection->query("select 'OK' AS result");

foreach ($result as $row) {
    echo "DB query result: ", $row['result'];
}
