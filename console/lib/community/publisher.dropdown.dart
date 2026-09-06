import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'api.dart' as api;
import 'publisher.typography.dart';
import 'publisher.list.row.dart';

class PublisherDropdown extends StatefulWidget {
  final ValueNotifier<api.PluginPublisher> current;
  final List<Widget> trailing;
  final List<Widget> leading;
  final Widget help;
  final bool readonly;
  final Future<api.PluginPublisherSearchResponse> Function(api.PluginPublisherSearchRequest) search;
  const PublisherDropdown({
    super.key,
    required this.current,
    this.trailing = const [],
    this.leading = const [],
    this.help = const ds.Hint(const Text("select from known social media platforms")),
    this.readonly = false,
    this.search = api.publishers.search,
  });

  @override
  State<PublisherDropdown> createState() => _PublisherDropdownState();
}

class _PublisherDropdownState extends State<PublisherDropdown> with ds.LoadingState {
  final TextEditingController _search = TextEditingController();
  // canRequestFocus is false so this button never leaves the FocusScope with
  // a focused descendant.
  final FocusNode _addFocus = FocusNode(canRequestFocus: false, skipTraversal: true);
  final ValueNotifier<int> _discovered = ValueNotifier<int>(0);
  Widget? _optional;

  void _refresh() {
    _search.clear();
    _search.text = widget.current.value.description;
  }

  @override
  void initState() {
    super.initState();
    widget.current.addListener(_refresh);
  }

  @override
  void dispose() {
    _addFocus.dispose();
    _discovered.dispose();
    super.dispose();
    widget.current.removeListener(_refresh);
  }

  @override
  Widget build(BuildContext context) {
    final defaults = ds.Defaults.of(context);

    return Column(
      mainAxisSize: MainAxisSize.min,
      spacing: defaults.spacing / 2,
      children: [
        ds.SearchDropdown.text(
          PublisherTypography.description(widget.current.value),
          padding: defaults.padding.copyWith(
            top: defaults.padding.left / 4,
            bottom: defaults.padding.right / 4,
          ),
          key: ValueKey(widget.current.value.id),
          controller: _search,
          textAlign: TextAlign.center,
          help: widget.help,
          refresh: _discovered,
          leading: widget.leading,
          trailing: widget.trailing,
          onSearch: (query, onClick) {
            return widget.search(api.PluginPublisherSearchRequest()..query = query).then((response) {
              if (response.items.length <= 1) return ds.Empty;
              return Container(
                constraints: BoxConstraints(maxHeight: 400),
                child: ListView.builder(
                  shrinkWrap: true,
                  itemCount: response.items.length,
                  itemBuilder: (context, index) {
                    final current = response.items[index];
                    if (current.id == widget.current.value.id) {
                      // no need to display the current item in the list.
                      return ds.Empty;
                    }

                    return PublisherRow(
                      current,
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
