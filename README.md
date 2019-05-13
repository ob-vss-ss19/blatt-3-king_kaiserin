## Ausführen mit Docker

-   Images bauen

    ```
    make docker
    ```

-   ein (Docker)-Netzwerk `actors` erzeugen

    ```
    docker network create actors
    ```

-   Starten des Tree-Services und binden an den Port 8090 des Containers mit dem DNS-Namen
    `treeservice` (entspricht dem Argument von `--name`) im Netzwerk `actors`:

    ```
    docker run --rm --net actors --name treeservice treeservice \
      --bind="treeservice.actors:8090"
    ```

    Damit das funktioniert, müssen Sie folgendes erst im Tree-Service implementieren:

    -   die `main` verarbeitet Kommandozeilenflags und
    -   der Remote-Actor nutzt den Wert des Flags
    -   wenn Sie einen anderen Port als `8090` benutzen wollen,
        müssen Sie das auch im Dockerfile ändern (`EXPOSE...`)

-   Starten des Tree-CLI, Binden an `treecli.actors:8091` und nutzen des Services unter
    dem Namen und Port `treeservice.actors:8090`:

    ```
    docker run --rm --net actors --name treecli treecli --bind="treecli.actors:8091" \
      --remote="treeservice.actors:8090" trees
    ```

    Hier sind wieder die beiden Flags `--bind` und `--remote` beliebig gewählt und
    in der Datei `treeservice/main.go` implementiert. `trees` ist ein weiteres
    Kommandozeilenargument, dass z.B. eine Liste aller Tree-Ids anzeigen soll.

    Zum Ausprobieren können Sie den Service dann laufen lassen. Das CLI soll ja jedes
    Mal nur einen Befehl abarbeiten und wird dann neu gestartet.

-   Zum Beenden, killen Sie einfach den Tree-Service-Container mit `Ctrl-C` und löschen
    Sie das Netzwerk mit

    ```
    docker network rm actors
    ```

## Ausführen mit Docker ohne vorher die Docker-Images zu bauen

Nach einem Commit baut der Jenkins, wenn alles durch gelaufen ist, die beiden
Docker-Images. Sie können diese dann mit `docker pull` herunter laden. Schauen Sie für die
genaue Bezeichnung in die Consolenausgabe des Jenkins-Jobs.

Wenn Sie die Imagenamen oben (`treeservice` und `treecli`) durch die Namen aus der
Registry ersetzen, können Sie Ihre Lösung mit den selben Kommandos wie oben beschrieben,
ausprobieren.


## CLI commands

-   Starten des Tree-CLI mit `go run main.go`
-   Tree-Service:
    -   Erstellen eines neuen Baums: `go run main.go --newTree`
    -   Einfügen eines Wertes: `go run main.go --insert --key=8 --value=acht --ID=1001 --token=abcde`
    -   Löschen eines Wertes: `go run main.go --delete --key=8 --ID=1001 --token=abcde`
    -   Traversieren: `go run main.go --traverse --ID=1001 --token=abcde`
    -   Suchen: `go run main.go --search --key=8 --ID=1001 --token=abcde`
    -   Baum löschen (muss man 2 Mal direkt hintereinander ausführen um endgültig zu löschen): `go run main.go --deleteTree --ID=1001 --token=abcde`
-  Wird der Token / die ID nicht gefunden, wird die entsprechende Fehlermeldung ausgegeben.

## Flags:
  - ID int
        ID of the Tree (default 1)
  - bind string
        Adresse to bind CLI (default "localhost:8091")
  - delete
        delete value and key from tree
  - deleteTree
        delete whole Tree
  - insert
        insert new value into the tree
  - key int
        Key which is needed for Insert/Search/Delete (default 1)
  - nameCLI string
        Name for the CLI (default "treecli")
  - nameService string
        Name for the Service (default "treeservice")
  - newTree
        creates new tree, prints out id and token
  - remote string
        Adresse to bind Service (default "localhost:8090")
  - search
        search value for a key
  - size int
        size of a leaf (default 1)
  - token string
        Token of the Tree
  - traverse
        go through tree and get sorted key-value-Pairs
  - value string
        Vale which is needed to insert new key-value-Pair
