import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'autoimport.api.dart' as api;
import 'autoimport.edit.dart';

String _debounce(ds.Int64 seconds) {
  final d = Duration(seconds: seconds.toInt());
  if (d.inMinutes < 60) return "${d.inMinutes}m";
  final minutes = d.inMinutes % 60;
  return minutes == 0 ? "${d.inHours}h" : "${d.inHours}h${minutes}m";
}

class AutoimportRow extends StatefulWidget {
  final api.AutoimportDirectory current;
  final void Function(api.AutoimportDirectory? upd) onChange;
  final api.FnAutoimportDelete delete;

  const AutoimportRow({
    super.key,
    required this.current,
    this.onChange = ds.fnNoop,
    this.delete = api.autoimport.delete,
  });

  @override
  State<AutoimportRow> createState() => _AutoimportRowState();
}

class _AutoimportRowState extends State<AutoimportRow> with ds.LoadingState {
  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return ds.ErrorScreen(
      ds.CompactingMenu(
        [
          ds.CompactingMenu.expanded(
            Text(widget.current.path, maxLines: 1, overflow: TextOverflow.ellipsis),
          ),
          if (widget.current.description.isNotEmpty)
            Text(widget.current.description, maxLines: 1, overflow: TextOverflow.ellipsis),
          ds.CompactingMenu.pinned(Text(api.autoimportModeLabel(widget.current.mode))),
          ds.CompactingMenu.pinned(Text(_debounce(widget.current.debounce))),
          ds.LoadingIconButton(
            icon: const Icon(Icons.delete_outline),
            tooltip: "stop monitoring directory",
            onPressed: () async {
              // the confirmation is mounted by the modal root, above the local daemon's
              // authorization cache, resolve the credentials from this row's context.
              final authz = authn.request(authn.AuthzCache.meta(context));
              ds.modals
                  .of(context)
                  ?.push(
                    ds.Confirmation.yesNo(
                      content: Text("Are you sure you want to stop monitoring ${widget.current.path}?"),
                      onConfirm: (context) {
                        httpx
                            .withRetry(() => widget.delete(widget.current.id, options: [authz]))
                            .then((resp) => widget.onChange(null))
                            .catchError((cause) => widget.onChange(null), test: httpx.ErrorsTest.err404)
                            .catchError((cause) {
                              setState(() {
                                this.cause = ds.Error.unknown(cause, onTap: reseterr);
                              });
                            })
                            .whenComplete(() {
                              ds.modals.of(context)?.push(null);
                            });
                      },
                      onCancel: (context) {
                        ds.modals.of(context)?.push(null);
                      },
                    ),
                  );
            },
          ),
        ],
        icon: const Icon(Icons.expand_more_rounded),
      ),
      cause: cause,
      tint: defaults.dangerTint,
      borderRadius: defaults.borderRadius,
    );
  }
}

class AutoimportItem extends StatefulWidget {
  final api.AutoimportDirectory current;
  final void Function(api.AutoimportDirectory? upd) onChange;
  final api.FnAutoimportUpdate update;
  final api.FnAutoimportDelete delete;

  const AutoimportItem({
    super.key,
    required this.current,
    this.onChange = ds.fnNoop,
    this.update = api.autoimport.update,
    this.delete = api.autoimport.delete,
  });

  @override
  State<AutoimportItem> createState() => _AutoimportItemState();
}

class _AutoimportItemState extends State<AutoimportItem> with ds.LoadingState {
  // edits are applied to a copy so an abandoned edit never leaks into the listing.
  late api.AutoimportDirectory _edited = widget.current.deepCopy();

  @override
  void didUpdateWidget(covariant AutoimportItem old) {
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

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.TableRow.single(
      AutoimportRow(current: widget.current, onChange: widget.onChange, delete: widget.delete),
      expanded: Container(
        padding: theme.buttonTheme.padding,
        child: ds.ErrorScreen(
          Column(
            mainAxisSize: MainAxisSize.min,
            spacing: defaults.spacing,
            children: [
              AutoimportEdit(current: _edited, onChange: (v) => setState(() => _edited = v)),
              ds.LoadingButton(const Text("save"), onPressed: save),
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
