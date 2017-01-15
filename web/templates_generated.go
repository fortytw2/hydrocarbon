// Package web is generated with ftmpl {{{v0.3.1}}}, do not edit!!!! */
package web

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"os"

	"github.com/fortytw2/hydrocarbon"
)

func init() {
	_ = fmt.Sprintf
	_ = errors.New
	_ = os.Stderr
	_ = html.EscapeString
}

// TMPLERRbase evaluates a template base.tmpl
func TMPLERRbase(title string, loggedIn bool, unread int) (string, error) {
	_template := "base.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>`)
	_w(fmt.Sprintf(`%s`, _escape(title)))
	_w(`</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="/hydrocarbon.min.css">

</head>
<body>
	<ul id="menu">
		<li><a href="/">Hydrocarbon</a></li>
`)
	if loggedIn {
		_w(`    <li>
        <a href="/feeds">Feeds `)
		_w(fmt.Sprintf(`%d`, unread))
		_w(`</a>
    </li>
`)
	} else {
		_w(`    <li class="right"><a href="/login">Login</a></li>
    <li class="right"><a href="/register">Register</a></li>
`)
	}
	_w(`	</ul>

	<div class="content">
`)
	_w(`	</div>

	<footer>
		(c) 2017 <a rel="nofollow" href="https://github.com/fortytw2/hydrocarbon">[GitHub]</a>[Twitter][Email]
	</footer>
</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLbase evaluates a template base.tmpl
func TMPLbase(title string, loggedIn bool, unread int) string {
	html, err := TMPLERRbase(title, loggedIn, unread)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template base.tmpl:" + err.Error())
	}
	return html
}

// TMPLERRfeed evaluates a template feed.tmpl
func TMPLERRfeed(title string, loggedIn bool, unread int, feed *hydrocarbon.Feed, posts []hydrocarbon.Post) (string, error) {
	_template := "feed.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
`)
	_w(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>`)
	_w(fmt.Sprintf(`%s`, _escape(title)))
	_w(`</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="/hydrocarbon.min.css">

</head>
<body>
	<ul id="menu">
		<li><a href="/">Hydrocarbon</a></li>
`)
	if loggedIn {
		_w(`    <li>
        <a href="/feeds">Feeds `)
		_w(fmt.Sprintf(`%d`, unread))
		_w(`</a>
    </li>
`)
	} else {
		_w(`    <li class="right"><a href="/login">Login</a></li>
    <li class="right"><a href="/register">Register</a></li>
`)
	}
	_w(`	</ul>

	<div class="content">
`)
	_w(`

<h1>`)
	_w(fmt.Sprintf(`%v`, feed.Name))
	_w(`</h1>

`)
	for _, post := range posts {
		_w(`	<h2> `)
		_w(fmt.Sprintf(`%v`, post.Title))
		_w(`</h2>
	<p> `)
		_w(fmt.Sprintf(`%v`, post.Content))
		_w(`</p>
	<br>
`)
	}
	_w(`	</div>

	<footer>
		(c) 2017 <a rel="nofollow" href="https://github.com/fortytw2/hydrocarbon">[GitHub]</a>[Twitter][Email]
	</footer>
</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLfeed evaluates a template feed.tmpl
func TMPLfeed(title string, loggedIn bool, unread int, feed *hydrocarbon.Feed, posts []hydrocarbon.Post) string {
	html, err := TMPLERRfeed(title, loggedIn, unread, feed, posts)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template feed.tmpl:" + err.Error())
	}
	return html
}

// TMPLERRfeeds evaluates a template feeds.tmpl
func TMPLERRfeeds(title string, loggedIn bool, unread int, feeds []hydrocarbon.Feed) (string, error) {
	_template := "feeds.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
`)
	_w(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>`)
	_w(fmt.Sprintf(`%s`, _escape(title)))
	_w(`</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="/hydrocarbon.min.css">

</head>
<body>
	<ul id="menu">
		<li><a href="/">Hydrocarbon</a></li>
`)
	if loggedIn {
		_w(`    <li>
        <a href="/feeds">Feeds `)
		_w(fmt.Sprintf(`%d`, unread))
		_w(`</a>
    </li>
`)
	} else {
		_w(`    <li class="right"><a href="/login">Login</a></li>
    <li class="right"><a href="/register">Register</a></li>
`)
	}
	_w(`	</ul>

	<div class="content">
`)
	_w(`
<h1>Feed Listing</h1>

`)
	for _, f := range feeds {
		_w(`	<h2> `)
		_w(fmt.Sprintf(`%v`, f.Name))
		_w(`</h2>
	<p> `)
		_w(fmt.Sprintf(`%v`, f.Description))
		_w(`</p>
	<a href="/feeds?id=`)
		_w(fmt.Sprintf(`%v`, f.ID))
		_w(`">link</a>
	<br>
`)
	}
	_w(`	</div>

	<footer>
		(c) 2017 <a rel="nofollow" href="https://github.com/fortytw2/hydrocarbon">[GitHub]</a>[Twitter][Email]
	</footer>
</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLfeeds evaluates a template feeds.tmpl
func TMPLfeeds(title string, loggedIn bool, unread int, feeds []hydrocarbon.Feed) string {
	html, err := TMPLERRfeeds(title, loggedIn, unread, feeds)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template feeds.tmpl:" + err.Error())
	}
	return html
}

// TMPLERRhome evaluates a template home.tmpl
func TMPLERRhome(title string, loggedIn bool, unread int) (string, error) {
	_template := "home.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>`)
	_w(fmt.Sprintf(`%s`, _escape(title)))
	_w(`</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="/hydrocarbon.min.css">

</head>
<body>
	<ul id="menu">
		<li><a href="/">Hydrocarbon</a></li>
`)
	if loggedIn {
		_w(`    <li>
        <a href="/feeds">Feeds `)
		_w(fmt.Sprintf(`%d`, unread))
		_w(`</a>
    </li>
`)
	} else {
		_w(`    <li class="right"><a href="/login">Login</a></li>
    <li class="right"><a href="/register">Register</a></li>
`)
	}
	_w(`	</ul>

	<div class="content">
`)
	_w(`
<h1>Welcome To Hydrocarbon</h1><br>

<p>
	Hydrocarbon is a no-nonsense, full text news reader, powered by
	plugins (supporting RSS, as well).
</p>
`)
	_w(`	</div>

	<footer>
		(c) 2017 <a rel="nofollow" href="https://github.com/fortytw2/hydrocarbon">[GitHub]</a>[Twitter][Email]
	</footer>
</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLhome evaluates a template home.tmpl
func TMPLhome(title string, loggedIn bool, unread int) string {
	html, err := TMPLERRhome(title, loggedIn, unread)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template home.tmpl:" + err.Error())
	}
	return html
}

// TMPLERRlogin evaluates a template login.tmpl
func TMPLERRlogin(title string, loggedIn bool, unread int) (string, error) {
	_template := "login.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>`)
	_w(fmt.Sprintf(`%s`, _escape(title)))
	_w(`</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="/hydrocarbon.min.css">

</head>
<body>
	<ul id="menu">
		<li><a href="/">Hydrocarbon</a></li>
`)
	if loggedIn {
		_w(`    <li>
        <a href="/feeds">Feeds `)
		_w(fmt.Sprintf(`%d`, unread))
		_w(`</a>
    </li>
`)
	} else {
		_w(`    <li class="right"><a href="/login">Login</a></li>
    <li class="right"><a href="/register">Register</a></li>
`)
	}
	_w(`	</ul>

	<div class="content">
`)
	_w(`
<div id="register">

<h1>Login</h1>

<form action="login" method="post">
  Email <input type="email" name="email"><br>
  Password <input type="password" name="pass"><br>
  <a href="password_reset">forgot password?</a><br>
  <input type="submit" value="Submit">
</form>

</div>
`)
	_w(`	</div>

	<footer>
		(c) 2017 <a rel="nofollow" href="https://github.com/fortytw2/hydrocarbon">[GitHub]</a>[Twitter][Email]
	</footer>
</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLlogin evaluates a template login.tmpl
func TMPLlogin(title string, loggedIn bool, unread int) string {
	html, err := TMPLERRlogin(title, loggedIn, unread)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template login.tmpl:" + err.Error())
	}
	return html
}

// TMPLERRpost evaluates a template post.tmpl
func TMPLERRpost(title string, loggedIn bool, unread int) (string, error) {
	_template := "post.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>`)
	_w(fmt.Sprintf(`%s`, _escape(title)))
	_w(`</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="/hydrocarbon.min.css">

</head>
<body>
	<ul id="menu">
		<li><a href="/">Hydrocarbon</a></li>
`)
	if loggedIn {
		_w(`    <li>
        <a href="/feeds">Feeds `)
		_w(fmt.Sprintf(`%d`, unread))
		_w(`</a>
    </li>
`)
	} else {
		_w(`    <li class="right"><a href="/login">Login</a></li>
    <li class="right"><a href="/register">Register</a></li>
`)
	}
	_w(`	</ul>

	<div class="content">
`)
	_w(`
<h1>view one post </h1>
`)
	_w(`	</div>

	<footer>
		(c) 2017 <a rel="nofollow" href="https://github.com/fortytw2/hydrocarbon">[GitHub]</a>[Twitter][Email]
	</footer>
</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLpost evaluates a template post.tmpl
func TMPLpost(title string, loggedIn bool, unread int) string {
	html, err := TMPLERRpost(title, loggedIn, unread)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template post.tmpl:" + err.Error())
	}
	return html
}

// TMPLERRregister evaluates a template register.tmpl
func TMPLERRregister(title string, loggedIn bool, unread int) (string, error) {
	_template := "register.tmpl"
	_escape := html.EscapeString
	var _ftmpl bytes.Buffer
	_w := func(str string) { _, _ = _ftmpl.WriteString(str) }
	_, _, _ = _template, _escape, _w

	_w(`
`)
	_w(`
<!doctype html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>`)
	_w(fmt.Sprintf(`%s`, _escape(title)))
	_w(`</title>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="/hydrocarbon.min.css">

</head>
<body>
	<ul id="menu">
		<li><a href="/">Hydrocarbon</a></li>
`)
	if loggedIn {
		_w(`    <li>
        <a href="/feeds">Feeds `)
		_w(fmt.Sprintf(`%d`, unread))
		_w(`</a>
    </li>
`)
	} else {
		_w(`    <li class="right"><a href="/login">Login</a></li>
    <li class="right"><a href="/register">Register</a></li>
`)
	}
	_w(`	</ul>

	<div class="content">
`)
	_w(`
<div id="register">

<h1>Registration</h1>

<form action="register" method="post">
  Email <input type="email" name="email"><br>
  Password <input type="password" name="password"><br>
  <a href="password_reset">already have an account? login</a><br>
  <input type="submit" value="Submit">
</form>

</div>
`)
	_w(`	</div>

	<footer>
		(c) 2017 <a rel="nofollow" href="https://github.com/fortytw2/hydrocarbon">[GitHub]</a>[Twitter][Email]
	</footer>
</body>
</html>
`)

	return _ftmpl.String(), nil
}

// TMPLregister evaluates a template register.tmpl
func TMPLregister(title string, loggedIn bool, unread int) string {
	html, err := TMPLERRregister(title, loggedIn, unread)
	if err != nil {
		_, _ = os.Stderr.WriteString("Error running template register.tmpl:" + err.Error())
	}
	return html
}
