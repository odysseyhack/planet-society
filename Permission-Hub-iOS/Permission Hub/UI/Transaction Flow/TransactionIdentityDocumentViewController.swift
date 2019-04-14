//
//  TransactionIdentityDocumentViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionIdentityDocumentViewController: PHTableViewController {

    // MARK: - Private properties

    private var transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(title: "", items: [
            .notification(
                type: .verification,
                text: "This company is verified"),
            .description(
                image: nil,
                date: transaction.date,
                title: "Identity document (passport)",
                description: "Please fill out your passport details."),
            .selectionDisclosure(text: "Country of issue"),
            .form(
                placeholder: "Passport number",
                text: "J12393496",
                keyboardType: .numbersAndPunctuation),
            .form(
                placeholder: "Expiration date",
                text: "02/2022",
                keyboardType: .numbersAndPunctuation)
        ])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        tableView.allowsSelection = false
    }
}
