<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>{{ .PageTitle }}</title>
  <style>
  label {
    width: 3em;
    display: inline-block;
  }
  button {
    width: 5em;
  }
  </style>
</head>
<body>
  <form>
    <label for="target">Target</label>
    <input placeholder="Insert Long URL ..." name="target" id="target" value="{{ .Target }}"/>
    <button type="reset">Reset</button>
    <br/>
    <label for="key">Key</label>
    <input placeholder="Insert Key ..." name="key" id="key" value="{{ .Key }}" />
    <button type="button" id="action">Add</button>
  </form>
</body>
</html>
