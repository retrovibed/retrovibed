import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'autoimport.api.dart' as api;
import 'autoimport.edit.dart';

class AutoImportRow extends StatelessWidget {
  final api.AutoimportDirectory current;

  const AutoImportRow({super.key, required this.current});

  @override
  Widget build(BuildContext context) {
    return ds.CompactingMenu(
      [
        ds.CompactingMenu.expanded(
          Text(
            current.description.isNotEmpty ? current.description : current.path,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ),
      ],
      icon: const Icon(Icons.expand_more_rounded),
    );
  }
}

class AutoImportItem extends StatefulWidget {
  final api.AutoimportDirectory current;
  final void Function(api.AutoimportDirectory? upd) onChange;
  final api.FnAutoimportUpdate update;
  final api.FnAutoimportDelete delete;

  const AutoImportItem({
    super.key,
    required this.current,
    this.onChange = ds.fnNoop,
    this.update = api.autoimport.update,
    this.delete = api.autoimport.delete,
  });

  @override
  State<AutoImportItem> createState() => _AutoImportItemState();
}

class _AutoImportItemState extends State<AutoImportItem> with ds.LoadingState {
  // edits are applied to a copy so an abandoned edit never leaks into the listing.
  late api.AutoimportDirectory _edited = widget.current.deepCopy();

  @override
  void didUpdateWidget(covariant AutoImportItem old) {
    super.didUpdateWidget(old);
    if (old.current != widget.current) _edited = widget.current.deepCopy();
  }

  Future<void> save() {
    return httpx
        .withRetry(
          () => widget.update(
            widget.current.id,
            api.AutoimportDirectoryUpdateRequest(directory: _edited),
            options: [authn.request(authn.AuthzCache.meta(context))],
          ),
        )
        .then((resp) => widget.onChange(resp.directory))
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });
  }

  Future<void> remove() {
    // the confirmation is mounted by the modal root, above the local daemon's authorization
    // cache, so the credentials are resolved from this item's context.
    final authz = authn.request(authn.AuthzCache.meta(context));
    return ds.modals
        .asyncfn(
          context,
          ds.Confirmation.dangerous(
            content: Text("Are you sure you want to stop monitoring ${widget.current.path}?"),
            onConfirm: (_) => httpx
                .withRetry(() => widget.delete(widget.current.id, options: [authz]))
                .then((_) => widget.onChange(null))
                .catchError((_) => widget.onChange(null), test: httpx.ErrorsTest.err404),
          ),
        )
        .catchError((cause) {
          setState(() {
            this.cause = ds.Error.unknown(cause, onTap: reseterr);
          });
        });
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.TableRow.single(
      AutoImportRow(current: widget.current),
      expanded: ds.Container(
        padding: defaults.padding,
        ds.ErrorScreen(
          AutoImportEdit(
            current: _edited,
            onChange: (v) => setState(() => _edited = v),
            actions: [
              ds.LoadingIconButton(
                icon: const Icon(Icons.save),
                tooltip: "save changes",
                onPressed: save,
              ),
              ds.LoadingIconButton(
                icon: const Icon(Icons.delete_outline),
                tooltip: "stop monitoring directory",
                onPressed: remove,
              ),
            ],
          ),
          cause: cause,
          tint: defaults.dangerTint,
          borderRadius: defaults.borderRadius,
        ),
      ),
    );
  }
}
