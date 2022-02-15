# go-api-books
REST API with unit tests and coverage percentage
API developed from a book CRUD, where unit tests are performed and their percentage of coverage

# Download Gin Framework
<pre><code>go get github.com/gin-gonic/gin</code></pre>

# Download GORM
<pre><code>go get github.com/jinzhu/gorm</code></pre>

# Download Driver Postgres
<pre><code>go get github.com/jinzhu/gorm/dialects/postgres</code></pre>

# Download GoDotEnv
<pre><code>go get github.com/joho/godotenv</code></pre>

# API
<h3>localhost:5000/</h3>
<ul>
<li><code>GET</code> : GET all book</li>
<li><code>POST</code> : POST create a new book</li>
</ul>

<h3>localhost:5000/:id</h3>
<ul>
<li><code>GET</code> : GET a book</li>
<li><code>PUT</code> : PUT update a book</li>
<li><code>DELETE</code> : DELETE a book</li>
</ul>

<h3>POST Params</h3>
<pre><code>{
	"author": "Op Super John Doe Bilw",
	"name": "Implementation Golang"
}
</code></pre>

<h3>PUT Params</h3>
<pre><code>
{
  "ID": 1,
	"author": "Op Super John Doe Bilw",
	"name": "Implementation Golang"
}
</code></pre>

# Test 
<h3>To create test files</h3>
<pre><code>gotests -all -w .</code></pre>

<h3>Create Tests</h3>
<pre><code>go test -v ./...</code></pre>

<h3>Create coverage</h3>
<ul>
  <h4>create cover profiler</h4>
  <li>
    <pre><code>go test -coverprofile<nameController></code></pre>
  </li>
  <h4>open cover profiler</h4>
  <li>
    <pre><code>go tool cover -html=nameCoverProfile</code></pre>
  </li>
</ul>
