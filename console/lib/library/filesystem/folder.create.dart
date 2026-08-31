import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/media.dart' as media;
import 'package:retrovibed/mimex.dart' as mimex;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;

// names a new folder inside parent. presented as a modal from the library menu, beside
// the other actions that add content, rather than as a permanent row.
class FolderCreate extends StatefulWidget {
  final String parent;
  final media.FnMkdir mkdir;
  final void Function(media.Media created) onCreated;
  final VoidCallback onCancel;

  const FolderCreate({
    super.key,
    required this.parent,
    required this.onCreated,
    required this.onCancel,
    this.mkdir = media.media.mkdir,
  });

  @override
  State<StatefulWidget> createState() => _FolderCreate();
}

class _FolderCreate extends State<FolderCreate> {
  final _name = TextEditingController();
  Widget _cause = ds.Error.zero;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void reseterr() {
    setState(() {
      _cause = ds.Error.zero;
    });
  }

  @override
  void dispose() {
    _name.dispose();
    super.dispose();
  }

  Future<void> submit() {
    final name = _name.text.trim();
    if (name.isEmpty) return Future.value();

    return httpx
        .withRetry(
          () => widget.mkdir(
            name,
            widget.parent,
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then((v) => widget.onCreated(v.media))
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unauthorized(cause, onTap: reseterr);
          });
        }, test: httpx.ErrorsTest.unauthorized)
        .catchError((cause) {
          setState(() {
            _cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    return ds.Card(
      leading: [ds.Heading(const Text("new folder"))],
      trailing: [
        IconButton(icon: const Icon(Icons.close), onPressed: widget.onCancel),
      ],
      forms.Container(
        forms.Field(
          cause: _cause,
          input: TextField(
            autofocus: true,
            controller: _name,
            decoration: InputDecoration(
              hintText: "folder name",
              icon: Icon(mimex.icofolder),
            ),
            onSubmitted: (_) => submit(),
          ),
          trailing: [
            ds.LoadingIconButton(onPressed: submit, icon: const Icon(Icons.check)),
          ],
        ),
      ),
    );
  }
}
