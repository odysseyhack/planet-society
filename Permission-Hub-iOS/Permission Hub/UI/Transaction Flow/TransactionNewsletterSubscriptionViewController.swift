//
//  TransactionNewsletterSubscriptionViewController.swift
//  Permission Hub
//
//  Created by Corné on 14/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionNewsletterSubscriptionViewController: PHTableViewController {

    // MARK: - Private properties

    private var transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(title: "", items: [
            .notification(
                type: .verification,
                text: "This company is verified"),
            .notification(
                type: .warning,
                text: "Permission warning!"),
            .description(
                date: transaction.date,
                title: "Newsletter subscription",
                description: "Phone House is using a newsletter to communicate with their clients and potential clients for marketing purposes."),
            .selection(options: [
                "Yes",
                "No"
            ])
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

