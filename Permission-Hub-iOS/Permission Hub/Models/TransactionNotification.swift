//
//  TransactionNotification.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import Foundation

struct TransactionItem: Decodable {
    let item: String
    let fields: [String]

    private enum CodingKeys: String, CodingKey {
        case item = "Item"
        case fields = "Fields"
    }
}

struct TransactionNotification: Decodable {
    let transactionID: String
    let items: [TransactionItem]
    let reason: String
    let verification: [String]
    let date: Date
    let requesterName: String
    let requesterPublicKey: String

    private enum CodingKeys: String, CodingKey {
        case transactionID
        case items = "item"
        case reason
        case verification
        case date
        case requesterName
        case requesterPublicKey = "RequesterPublicKey"
    }

    // MARK: - Initialization

    init(
        transactionID: String,
        items: [TransactionItem],
        reason: String,
        verification: [String],
        date: Date,
        requesterName: String,
        requesterPublicKey: String) {

        self.transactionID = transactionID
        self.items = items
        self.reason = reason
        self.verification = verification
        self.date = date
        self.requesterName = requesterName
        self.requesterPublicKey = requesterPublicKey
    }

    init(from decoder: Decoder) throws {

        let container = try decoder.container(keyedBy: CodingKeys.self)
        let transactionID = try container.decode(String.self, forKey: .transactionID)
        let items = try container.decode([TransactionItem].self, forKey: .items)
        let reason = try container.decode(String.self, forKey: .reason)
        let verification = try container.decode([String].self, forKey: .verification)

        let dateString = try container.decode(String.self, forKey: .date)
        let dateFormatter = ISO8601DateFormatter()

        guard let date = dateFormatter.date(from: dateString) else {
            throw DecodingError.valueNotFound(
                Date.self,
                DecodingError.Context(
                    codingPath: [],
                    debugDescription: "Error decoding date object"))
        }

        let requesterName = try container.decode(String.self, forKey: .requesterName)
        let requesterPublicKey = try container.decode(String.self, forKey: .requesterPublicKey)

        self.init(
            transactionID: transactionID,
            items: items,
            reason: reason,
            verification: verification,
            date: date,
            requesterName: requesterName,
            requesterPublicKey: requesterPublicKey)
    }
}
