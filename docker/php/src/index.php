<h1>Welcome to <?php echo getenv("APP_URL"); ?></h1><br/>

<p>Here is the list of tables in your database:</p> <br/>
<?php
$servername = getenv("DB_HOST");
$username = getenv("DB_USER");
$password = getenv("DB_PASSWORD");
$dbname = getenv("DB_NAME");

// Create connection
$conn = new mysqli($servername, $username, $password, $dbname);
// Check connection
if ($conn->connect_error) {
  die("Connection failed: " . $conn->connect_error);
}

$sql = "show tables;";
$result = $conn->query($sql);

if ($result->num_rows > 0) {
  // output data of each row
  while($row = $result->fetch_assoc()) {
    echo $row["Tables_in_$dbname"] . "<br>";
  }
} else {
  echo "0 results";
}
$conn->close();
?> 