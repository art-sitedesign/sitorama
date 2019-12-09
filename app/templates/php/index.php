<?php

{{range $i, $checker := .Checkers}}
{{$checker}}
echo '<hr>';

{{end}}

?>