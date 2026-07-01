import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as api;
import 'daemon.auto.dart';
import 'daemon.item.dart';
import 'daemon.manual.dart';
import 'daemon.typography.dart';

class DaemonDropdown extends StatefulWidget {
  final ValueNotifier<api.Daemon> library;
  final List<Widget> trailing;
  final Widget help;
  const DaemonDropdown({
    super.key,
    required this.library,
    this.trailing = const [],
    this.help = const ds.Hint(const Text("select which daemon instance to configure from the dropdown")),
  });

  @override
  State<DaemonDropdown> createState() => _DaemonDropdownState();
}

class _DaemonDropdownState extends State<DaemonDropdown> {
  final TextEditingController _search = TextEditingController();
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
  }

  @override
  void dispose() {
    super.dispose();
    widget.library.removeListener(_refresh);
  }

  @override
  Widget build(BuildContext context) {
    return ds.Help(
      Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ds.SearchDropdown.text(
            DaemonTypography.description(widget.library.value),
            key: ValueKey(widget.library.value.id),
            controller: _search,
            textAlign: TextAlign.center,
            leading: [
              IconButton(
                tooltip: "connect to another library",
                onPressed: () {
                  setState(() {
                    _optional = _optional != null
                        ? null
                        : ManualConfiguration(
                            connect: (daemon) {
                              setState(() {
                                _optional = null;
                              });
                              EndpointAuto.of(
                                context,
                              )?.setdaemon(daemon).ignore();
                            },
                          );
                  });
                },
                icon: Icon(_optional == null ? Icons.add : Icons.remove),
              ),
            ],
            trailing: widget.trailing,
            onSearch: (query, onClick) {
              return api.daemons.search(api.DaemonSearchRequest()..query = query).then((response) {
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

                      return DaemonDropdownItem(
                        library: daemon,
                        onTap: onClick,
                      );
                    },
                  ),
                );
              });
            },
          ),
          ?_optional,
        ],
      ),
      widget.help,
    );
  }
}
