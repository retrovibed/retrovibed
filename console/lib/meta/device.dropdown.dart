import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'api.dart' as api;
import 'device.auto.dart';
import 'device.item.dart';
import 'device.manual.dart';
import 'device.typography.dart';

class DaemonDropdown extends StatefulWidget {
  final ValueNotifier<api.Daemon> library;
  final List<Widget> trailing;
  final List<Widget> leading;
  final Widget help;
  final bool readonly;
  final bool remoteonly;
  final Future<api.DaemonSearchResponse> Function(api.DaemonSearchRequest) search;
  final Future<Stream<api.Daemon>> Function({List<httpx.Option> options}) discover;
  final DaemonOnSelect onSelect;
  const DaemonDropdown({
    super.key,
    required this.library,
    this.trailing = const [],
    this.leading = const [],
    this.help = const ds.Hint(const Text("select which daemon instance to configure from the dropdown")),
    this.remoteonly = false,
    this.readonly = false,
    this.search = api.daemons.search,
    this.discover = api.daemons.discover,
    this.onSelect = global,
  });

  // validates via EndpointAuto and mutates the app-wide active host (httpx.set + EndpointAuto.changed).
  static Future<api.Daemon> global(BuildContext context, api.Daemon daemon) {
    return (EndpointAuto.of(context)?.refreshNoErrHandling(Future.value(daemon)) ?? Future.value()).then((_) => daemon);
  }

  // validates connectability only; never touches EndpointAuto/httpx globals.
  static Future<api.Daemon> local(BuildContext context, api.Daemon daemon) {
    return api.daemons.connectable(daemon);
  }

  @override
  State<DaemonDropdown> createState() => _DaemonDropdownState();
}

class _DaemonDropdownState extends State<DaemonDropdown> {
  final TextEditingController _search = TextEditingController();
  // canRequestFocus is false so this button never leaves the FocusScope with
  // a focused descendant, which would otherwise block ManualConfiguration's
  // autofocus when it opens.
  final FocusNode _addFocus = FocusNode(canRequestFocus: false, skipTraversal: true);
  final ValueNotifier<int> _discovered = ValueNotifier<int>(0);
  bool _scanning = false;
  Widget? _optional;

  void setState(VoidCallback fn) {
    if (!mounted) return;
    super.setState(fn);
  }

  void _refresh() {
    _search.clear();
    _search.text = widget.library.value.description;
  }

  @override
  void initState() {
    super.initState();
    widget.library.addListener(_refresh);
    ds.postframe(_scanForPeers);
  }

  Future<void> _scanForPeers() async {
    setState(() => _scanning = true);
    widget
        .discover()
        .then((s) {
          return s.forEach((d) {
            _discovered.value++;
          });
        })
        .whenComplete(() => setState(() => _scanning = false))
        .catchError((cause) {
          debugPrint(cause.toString());
        });
  }

  @override
  void dispose() {
    _addFocus.dispose();
    _discovered.dispose();
    super.dispose();
    widget.library.removeListener(_refresh);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);
    return Column(
      mainAxisSize: MainAxisSize.min,
      spacing: defaults.spacing / 2,
      children: [
        ds.SearchDropdown.text(
          DaemonTypography.description(widget.library.value),
          padding: defaults.padding.copyWith(
            top: defaults.padding.left / 4,
            bottom: defaults.padding.right / 4,
          ),
          key: ValueKey(widget.library.value.id),
          controller: _search,
          textAlign: TextAlign.center,
          help: widget.help,
          refresh: _discovered,
          leading: [
            ...widget.leading,
            if (!widget.readonly)
              ds.LoadingIconButton(
                tooltip: "connect to another library",
                focusNode: _addFocus,
                onPressed: () async {
                  setState(() {
                    _optional = _optional != null
                        ? null
                        : ds.Container(
                            padding: defaults.padding.copyWith(top: 0),
                            ManualConfiguration(
                              connect: (daemon) {
                                setState(() {
                                  _optional = null;
                                });
                                EndpointAuto.of(
                                  context,
                                )?.setdaemon(daemon).ignore();
                              },
                            ),
                          );
                  });
                },
                icon: Icon(_optional == null ? Icons.add : Icons.remove),
              ),
          ],
          trailing: [
            if (_scanning)
              Padding(
                padding: const EdgeInsets.all(4),
                child: SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(strokeWidth: 2.0),
                ),
              ),
            ...widget.trailing,
          ],
          onSearch: (query, onClick) {
            return widget.search(api.DaemonSearchRequest()..query = query).then((response) {
              if (response.items.length <= 1) return ds.Empty;
              return Container(
                constraints: BoxConstraints(maxHeight: 400),
                child: ListView.builder(
                  shrinkWrap: true,
                  itemCount: response.items.length,
                  itemBuilder: (context, index) {
                    final daemon = response.items[index];
                    if (daemon.id == widget.library.value.id) {
                      // no need to display the current library in the list.
                      return ds.Empty;
                    }
                    if (widget.remoteonly && api.daemons.isLocalDevice(daemon)) {
                      return ds.Empty;
                    }

                    return DaemonDropdownItem(
                      library: daemon,
                      readonly: widget.readonly,
                      onSelect: widget.onSelect,
                      onTap: (v) {
                        widget.library.value = v;
                        onClick();
                      },
                    );
                  },
                ),
              );
            });
          },
        ),
        ?_optional,
      ],
    );
  }
}
