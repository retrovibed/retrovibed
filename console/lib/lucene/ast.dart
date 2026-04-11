// Lucene query AST.
//
// Represents parsed structure without field-level semantics.
// Field types (Boolean, Timestamp, etc.) interpret nodes during apply().

enum Occur { must, should, mustNot }

class Query {
  static final Query empty = Query(const []);

  final List<Clause> clauses;
  const Query(this.clauses);
}

class Clause {
  final Occur occur;
  final Node node;

  const Clause(this.occur, this.node);

  factory Clause.must(Node n) => Clause(Occur.must, n);
  factory Clause.should(Node n) => Clause(Occur.should, n);
  factory Clause.mustNot(Node n) => Clause(Occur.mustNot, n);
}

sealed class Node {}

// field:value — field is null for bare words/phrases
class Term extends Node {
  final String field;
  final String value;
  Term(this.field, this.value);
}

// field:[min TO max] or field:{min TO max}
// comparison operators (>x, >=x, <x, <=x) are normalised to Range during parsing
class Range extends Node {
  final String field;
  final String? min; // null means *
  final String? max; // null means *
  final bool minInclusive; // true for [, false for {
  final bool maxInclusive; // true for ], false for }

  Range(
    this.field, {
    this.min,
    this.max,
    this.minInclusive = true,
    this.maxInclusive = true,
  });
}

// (sub-query)
class Group extends Node {
  final Query query;
  Group(this.query);
}
