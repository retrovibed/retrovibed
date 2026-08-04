import 'package:flutter/material.dart';
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/design.kit/forms.dart' as forms;
import 'package:retrovibed/ddisc.dart' as ddisc;
import 'package:retrovibed/library.dart' as lib;

class DiscoveryDetails extends StatelessWidget {
  final ddisc.Discovery current;
  final lib.Known known;
  final EdgeInsets margin;
  final Widget help;

  const DiscoveryDetails(
    this.current,
    this.known, {
    super.key,
    this.margin = EdgeInsets.zero,
    this.help = ds.HelpScope.None,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final defaults = ds.Defaults.of(context);

    return ds.Card(
      alignment: Alignment.topLeft,
      margin: margin,
      help: help,
      SingleChildScrollView(
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          mainAxisSize: MainAxisSize.min,
          spacing: defaults.spacing / 4,
          children: [
            Text(
              known.description.isNotEmpty ? known.description : current.title,
              style: theme.textTheme.titleMedium,
            ),
            Text("Media", style: theme.textTheme.titleSmall),
            forms.Field(
              label: const Text("summary"),
              input: Text(known.summary.isEmpty ? '—' : known.summary),
            ),
            forms.Field(
              label: const Text("rating"),
              input: ds.Rating(rating: known.rating),
            ),
            forms.Field(
              label: const Text("released"),
              input: ds.Timestamp.iso8601(
                known.released,
                neginf: Text("unknown"),
              ),
            ),
            forms.Field(
              label: const Text("mimetype"),
              input: Text(known.mimetype.isEmpty ? '—' : known.mimetype),
            ),
            forms.Field(
              label: const Text("source"),
              input: Text(known.source.isEmpty ? '—' : known.source),
            ),
            forms.Field(label: const Text("adult"), input: Text(known.adult ? "yes" : "no")),
            Text("Discovery", style: theme.textTheme.titleSmall),
            forms.Field(label: const Text("size"), input: ds.Bytes(current.bytes)),
            forms.Field(label: const Text("health"), input: Text("${current.health}")),
            forms.Field(label: const Text("attempts"), input: Text("${current.attempts}")),
            forms.Field(
              label: const Text("discovered"),
              input: ds.Timestamp.iso8601(current.createdAt, neginf: ds.Empty),
            ),
            forms.Field(
              label: const Text("updated"),
              input: ds.Debug.pink(ds.Timestamp.iso8601(current.updatedAt, neginf: ds.Empty)),
            ),
            forms.Field(
              label: const Text("next check"),
              input: ds.Timestamp.iso8601(current.nextCheck, neginf: ds.Empty),
            ),
            forms.Field(label: const Text("updated"), input: Text(current.updatedAt)),
            forms.Field(label: const Text("poster"), input: Text(known.image)),
            forms.Field(
              label: const Text("discovered"),
              input: Text(current.source.isEmpty ? '—' : current.source),
            ),
            forms.Field(label: const Text("id"), input: Text(current.id)),
          ],
        ),
      ),
    );
  }
}
